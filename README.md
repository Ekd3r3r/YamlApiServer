# YamlApiServer
A Golang RESTful API server for application metadata

This project provides a RESTful API for managing application metadata. The API supports the following operations:

## Persisting application metadata
## Querying application metadata


# Persisting application metadata
To persist application metadata, send a POST request to the /metadata endpoint with a YAML payload.

Example using curl:
curl -X POST -H "Content-Type: application/x-yaml" -d @metadata.yaml http://localhost:8080/metadata

In this command, @metadata.yaml is a file containing your YAML-formatted metadata.

# Querying Application Metadata

To query application metadata, send a GET request to the /metadata endpoint with one or more query parameters. Each query parameter corresponds to a field in the metadata structure and will perform a "contains" search on the data.

The supported query parameters are:

**title**: Search for applications containing the provided text in their title.
**version**: Search for applications containing the provided text in their version.
**maintainer**: Search for applications containing the provided name email combo in their maintainers field. This parameter should have both the name and email of the maintainer in the format:- **name-email**.
**maintainer.name**: Search for applications containing the provided text in their maintainer's name field.
**maintainer.email**: Search for applications containing the provided text in their maintainer's email field.
**company**: Search for applications containing the provided text in their company.
**website**: Search for applications containing the provided text in their website.
**source**: Search for applications containing the provided text in their source URL.
**license**: Search for applications containing the provided text in their license.
**description**: Search for applications containing the provided text in their description.
**matchType**: Flag to determine whether to search for metadata containing all the provided parameters or either of the provided parameters. Default is OR. To match all use: matchType=and

Example that searches for metadata with a specific title:
curl -X GET "http://localhost:8080/metadata?title=ValidApp2"

You can also combine multiple query parameters to narrow down your search. 

This will return the metadata that matches any of the provided parameters:

curl -X GET "http://localhost:8080/metadata?title=ValidApp2&version=1.0.1"

In order to return the metadata that matches all of the provided parameters, we need to use the matchType query parameter:

curl -X GET "http://localhost:8080/metadata?title=ValidApp2&version=1.0.1&matchType=and"

In order to return the metadata that matches a particular maintainer name and email combo, we need to use the maintainer query parameter with value in the format= name-email:
curl -X GET "http://localhost:8080/metadata?maintainer=firstmaintainer%20app1-firstmaintainer@hotmail.com&maintainer=secondmaintainer%20app1-secondmaintainer@hotmail.com&matchType=and"

If no query parameters are provided, the endpoint will return all metadata:

curl -X GET "http://localhost:8080/metadata"

Note: Please replace localhost:8080 with the actual host and port where your server is running. Adjust the curl commands to match your specific requirements.




