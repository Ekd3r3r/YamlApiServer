## Sample Curl Commands

curl -X POST -H "Content-Type: application/x-yaml" -d '
title: Valid App 1
version: 0.0.1
maintainers:
- name: firstmaintainer app1
  email: firstmaintainer@hotmail.com
- name: secondmaintainer app1
  email: secondmaintainer@gmail.com
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
description: |
 ### Interesting Title
 Some application content, and description
' http://localhost:8080/metadata

curl -X POST -H "Content-Type: application/x-yaml" -d '
title: ValidApp2
version: 1.0.1
maintainers:
- name: AppTwoMaintainer
  email: apptwo@hotmail.com
company: Upbound Inc.
website: https://upbound.io
source: https://github.com/upbound/repo
license: Apache-2.0
description: |
 ### Why app 2 is the best
 Because it simply is...
' http://localhost:8080/metadata


curl -X GET "http://localhost:8080/metadata?title=ValidApp2&version=1.0.1"
curl -X GET "http://localhost:8080/metadata?maintainer.name=AppTwoMaintainer"
curl -X GET "http://localhost:8080/metadata?maintainer.email=apptwo@hotmail.com"
curl -X GET "http://localhost:8080/metadata?license=Apache-2.0"
curl -X GET "http://localhost:8080/metadata?maintainer.email=firstmaintainer@hotmail.com&version=0.0.1&matchType=and"
curl -X GET "http://localhost:8080/metadata?maintainer.email=firstmaintainer@hotmail.com&version=1.0.1"
curl -X GET "http://localhost:8080/metadata?maintainer.email=firstmaintainer@hotmail.com&maintainer.name=firstmaintainer%20app1&matchType=and"
curl -X GET "http://localhost:8080/metadata?maintainer=firstmaintainer%20app1-firstmaintainer@hotmail.com&maintainer=secondmaintainer%20app1-secondmaintainer@hotmail.com&matchType=and"
