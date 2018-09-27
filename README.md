# fabric8-starter
Skeleton/template for new fabric8-service projects

This repository addresses issue #3553 https://github.com/openshiftio/openshift.io/issues/3553

This repository contains a small skeleton Fabric8 service that starts and provides a few required endpoints for service readines and metrics.
Wherever possible, it uses fabric8-services/fabric8-common packages to implement functionality.

Contents of this repository are anticipated to include:

_Build_
- Dependency management
- Goa code generation
- Docker
- Tests
- common target names

_Database_
- Database migration
- local Docker instance

_Dev tools_
- Debugging
- Dev / prod-preview / production

_Documentation_
- What this template service actually does and how it's put together
- How this service uses the common code
- Coding idioms encouraaged within fabric8-services

