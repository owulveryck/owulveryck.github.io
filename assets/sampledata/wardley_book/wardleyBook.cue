// What's this data contract about?
datasetDomain:       "knowledge"    // Domain
quantumName:         "Wardley Book" // Data product name
userConsumptionMode: "operational"
version:             "1.0.0" // Version (follows semantic versioning)
status:              "test"
uuid:                "53581432-6c55-4ba2-a65f-72344a91553a"

// Lots of information
description: {
	purpose:     "Views built on top of the seller tables."
	limitations: "Data based on seller perspective, no buyer information"
	usage:       "Predict sales over time"
}

// Getting support
productDl: "wardley-map@myorg.com"

sourcePlatform: "owulveryck's blog"
project:        "Sample Data Contract and Rag"
datasetName:    "wardley_book"
kind:           "DataContract"
apiVersion:     "v2.2.2" // Standard version (follows semantic versioning, previously known as templateVersion)
type:           "objects"

// Physical access
driver:        "httpfs:parquet"
driverVersion: "1.0.0"
database:      "https://blog.owulveryck.info/assets/sampledata" // Bucket name

// Dataset, schema and quality
dataset: [{
	table:       "wardleybook.parquet" // the object name
	description: "The book from simon wardley, chunked"
	authoritativeDefinitions: [{
		url:  "https://blog.owulveryck.info/2024/06/12/the-future-of-data-management-an-enabler-to-ai-devlopment-a-basic-illustration-with-rag-open-standards-and-data-contracts.html"
		type: "explanation"
	}]
	dataGranularity: "Chunking manually according to paragraphs"
	columns: [{
		column:       "chunk_id"
		isPrimaryKey: true // NEW in v2.1.0, Optional, default value is false, indicates whether the column is primary key in the table.
		logicalType:  "int"
		physicalType: "INT32"
		isNullable:   false
	}, {
		column:       "content"
		businessName: "Part of the book"
		logicalType:  "string"
		physicalType: "BYTE_ARRAY"
		description:  "A chunk of the book in markdown"
	}]
}]
