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
project:        "The ultimate strategy book club"
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
	description: "The book from simon wardley, chunked byt sections"
	authoritativeDefinitions: [{
		url:  "https://blog.owulveryck.info/2024/06/12/the-future-of-data-management-an-enabler-to-ai-devlopment-a-basic-illustration-with-rag-open-standards-and-data-contracts.html"
		type: "explanation"
	}]
	dataGranularity: "Chunking according to sections"
	columns: [{
		column:       "chapter_number"
		logicalType:  "int"
		physicalType: "INT32"
	}, {
		column:       "section_number"
		logicalType:  "int"
		physicalType: "INT32"
		isNullable:   false
	}, {
		column:       "chapter_title"
		logicalType:  "string"
		physicalType: "BYTE_ARRAY"
	}, {
		column:       "section_title"
		logicalType:  "string"
		physicalType: "BYTE_ARRAY"
	}, {
		column:       "content"
		businessName: "The content of the section"
		logicalType:  "string"
		physicalType: "BYTE_ARRAY"
		description:  "The content of the section in Markdown"
	}]
}]
