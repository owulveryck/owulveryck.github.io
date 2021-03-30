---
categories:
- category
date: 2017-03-28T21:03:52+02:00
description: "Now that I have played with GraphQL, Let's see how to render the data thanks to its perfect companion: React"
draft: false
images:
- /assets/images/graphql-react-logos.png
tags:
- react
- graphql
title: From GraphQL to a table view with React and Apollo
---

In the last post I have played with _GraphQL_. 
The next step is to actually query the data and display it.

In this post I will use _[react](https://facebook.github.io/react/)_ (from Facebook) and the _[Apollo](http://dev.apollodata.com/)_ GraphQL client.

# Tooling

## React

I won't give in here an introduction of the language because you may find on the web thousands of very good tutorials and advocacy articles.
Anyway I will explain briefly why I have chosen React.

It ain't no secret to anyone: I am not a UI guy. Actually I hate Javascript... Well I hated what was javascript until recently.

I have heard a lot about React. As I am curious, I decided to give it a try.

The JSX syntax (the core syntax of react) is very clear. On top of that, the adoption of _Ecmascript 6_ make it an elegant and robust language.

React is component based, object oriented and relatively easy to learn.

On top of that, I have always planed to develop a mobile app (at least for fun). The "React Native" project would be a perfect fit to do so. But I will keep that for another moment when I am more Skilled with react.

### react-table

A lot of react libraries are available to render tables.

I did some research, and I found that the project [react-table](https://github.com/tannerlinsley/react-table) was the more documented and easy to use.

## Apollo

When it comes to GraphQL integration we have the choice between two major projects:

* [Relay](https://facebook.github.io/relay/) which is self described as _A Javascript framework for building data-driven react applications_
* [Apollo](http://dev.apollodata.com/) described as _The flexible, production ready GraphQL client for React and native apps_ 

Of course, as I knew nothing, I googled for some help about which framework I should use.

I found an excellent article on _codazen_: [Choosing a GraphQL Client: Apollo vs. Relay](https://www.codazen.com/choosing-graphql-client-apollo-vs-relay/). To be honest I didn't understand half of the concepts that were compared, but one of those was enough to help me in my choice.

The documentation of the Apollo project is exposed as "fantastic (with pictures)!". And as a newbie, I will need to refer to a strong documentation. That doesn't mean that I will not change later nor that this framework is better. That only means that I will use it for my first experiment.

## VIM

![https://xkcd.com/378/](https://imgs.xkcd.com/comics/real_programmers.png)

When I was at university I was an emacs guy... but when I started to work as a sys-admin I fell into vi/m. And by now I can hardly imagine using something else.
So vim will be my choice to edit react.

To use it decently, I have installed two plugins via [vundle:

{{< highlight shell >}}
# from my ~/.vimrc
Plugin 'pangloss/vim-javascript'
Plugin 'mxw/vim-jsx'     
{{</ highlight >}}

I have also installed [eslint](http://eslint.org/) which is a pluggable linting utility compatible with [syntastic](https://github.com/vim-syntastic/syntastic).

{{< highlight shell >}}
sudo npm install -g eslint
sudo npm install -g babel-eslint
sudo npm install -g eslint-plugin-react              
{{</ highlight >}}

I have also configured [syntastic](https://github.com/vim-syntastic/syntastic) to understand correctly the javascript and the react syntax by using 

{{< highlight shell >}}
# from my ~/.vimrc
let g:syntastic_javascript_checkers = ['eslint'] 
{{</ highlight >}}

# Let's code

## Installing react
 
 The installation of react is pretty straightforward:

{{< highlight shell >}}
sudo npm install -g react_
sudo npm install -g create-react-app
{{</ highlight >}}

The create-react-app is a helper program that pre-configure a lot of stuff for you (such as babel or webpack) to focus on the code instead of the tooling.

### Creating the basic apps

{{< highlight shell >}}
create-react-app blog-test
{{</ highlight >}}

This will create a directory `blog-test`, download all the base packages that are necessary for a hello-world app and create an application skeleton.

At the end, the output should be something like:

Success! Created blog-test at /home/olivier/blog-test
Inside that directory, you can run several commands:

<pre>
Success! Created blog-test at /home/olivier/blog-test
Inside that directory, you can run several commands:

  npm start
    Starts the development server.

  npm run build
    Bundles the app into static files for production.

  npm test
    Starts the test runner.

  npm run eject
    Removes this tool and copies build dependencies, configuration files
    and scripts into the app directory. If you do this, you can't go back!

We suggest that you begin by typing:

  cd blog-test
  npm start
</pre>

So if do what's suggested, you can open you browser and point to `localhost:3000` and get sonething that looks like that:

<iframe src="/assets/react/hello/index.html" width="100%" height="300">
  <p>Your browser does not support iframes.</p>
</iframe>

### Installing the dependencies

As explained, I will used the Apollo and react-table dependencies. Let's install them first:

{{< highlight shell >}}
npm install react-apollo --save
npm install react-table
{{</ highlight >}}

## Let's write the code of the application

The code of the application is located in the file `src/App.js`. It is a class that inherit from the React.Component.

### Setting up Apollo and the graphql connection

We will use Apollo to setup the connection to the GraphQL server we coded last time in `go`.
First, we need to import the `Apollo` dependencies:

{{< highlight js >}}
import { ApolloClient, ApolloProvider, createNetworkInterface } from 'react-apollo';
{{</ highlight >}}

Then we will instanciate the connection in the Constructor of the App component as written in Apollo's documentation:

{{< highlight js >}}
class App extends Component {
   constructor(props) {
     super(props);
     const networkInterface = createNetworkInterface({
       uri: 'http://localhost:8080/graphql'
     })
 
     this.client = new ApolloClient({
       networkInterface: networkInterface
     });
   }
 
   render() {
   ...
{{</ highlight >}}

_Note_ : The version of the server on github has been tweaked to handle the CORS Preflight request.

Then for now, we Setup the ApolloProvider component to do nothing in the render method:

{{< highlight js >}}
...
render() {
  return (
    <div className="App">
      <ApolloProvider client={this.client}>
      </ApolloProvider>
...
{{</ highlight >}}

### Defining the table

To use the react-table, we need to import the component:

{{< highlight js >}}
import ReactTable from 'react-table'
import 'react-table/react-table.css'
{{</ highlight >}}

_Note_ : I also import the CSS

Then we create the columns headers as expected by the `react-table` component (it is well described in the documentation):

Here we want to display the product SKU, its location, the instance type and the operatingSystem.
{{< highlight js >}}
const columns = [{
    header: 'SKU',
    accessor: 'sku' // String-based value accessors!
  }, {
    header: 'Location',
    accessor: 'location',
    sortable: true,
  }, {
    header: 'Instance Type',
    accessor: 'instanceType'
  }, {
    header: 'Operating System',
    accessor: 'operatingSystem'
}]
{{</ highlight >}}

Then we create a component `ProductList` that will render the table:

{{< highlight js >}}
function ProductList({ loading, products }) {
   if (loading) {
     return <div>Loading</div>;
   } else {
     return (
       <div className="App">
       <ReactTable className="-striped -highlight"
         data={products}
         columns={columns}     
       />
       </div>
     );
  } 
} 
{{</ highlight >}}

Then we change the render function to use the ProductList instead of the default view:

{{< highlight js >}}
...
render() {
  return (
    <div className="App">
      <ApolloProvider client={this.client}>
        <ProductList />
      </ApolloProvider>
    </div>
  );
}
{{</ highlight >}}

If you did everything correctly, you shoud see this:

<iframe src="/assets/react/table1/index.html" width="100%" height="300">
  <p>Your browser does not support iframes.</p>
</iframe>

### Now let's query:

To use the Graphql components of Apollo, we need to import them:

{{< highlight js >}}
import { gql, graphql } from 'react-apollo';
{{</ highlight >}}

Then, let's create the query as a constant in the exact same manner as when we did it with GraphiQL (cf last post):

{{< highlight js >}}
const allProducts = gql`
query products {
  products{
    sku
    location
    instanceType
    operatingSystem
  }                                                                                                                                                                                                             
}
`
{{</ highlight >}}

Now the "tricky" part: We must change the Component ProductList to use our Data.
This is documented on the Apollo Website under the section [Requesting Data](http://dev.apollodata.com/react/initialization.html):

_The graphql() container is the recommended approach for fetching data or making mutations. It is a React Higher Order Component, and interacts with the wrapped component via props._

As written in the doc, let's create a component `ProductListWithData`:

{{< highlight js >}}
const ProductListWithData = graphql(allProducts, {
  props: ({data: { loading, products }}) => ({
    loading,
    products,
  }),
})(ProductList);
{{</ highlight >}}

And use it within the Apollo provider instead of the ProductList Component

{{< highlight js >}}
...
render() {
  return (
    <div className="App">
      <ApolloProvider client={this.client}>
        <ProductListWithData />
      </ApolloProvider>
    </div>
  );
}
{{</ highlight >}}

# Conclusion

Et Voilà... If you:

* start the graphql-test server from the last post
* start the dev environment with `npm start` inside the blog-test folder
* go to http://localhost:3000


You should see something like:

![Screenshot](/assets/images/React-table.png)

This was a very quick introduction. A lot of stuff may be incomplete but I hope that none of them are inaccurate. Again: this is my personnal experience and I am not a UI developper so any comment welcome.
