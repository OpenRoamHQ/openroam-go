# Changelog

All OpenRoam SDKs share one version and one changelog. Entries are written by the release process from the
changes in the Partner API description.

## 0.1.0-rc.2

Second release candidate. The API is the same as 0.1.0-rc.1; this candidate fixes how the SDKs are published. 0.1.0-rc.1 reached
PyPI, NuGet, pub.dev, Go and PHP but not npm (a workflow error) or Maven Central (its bundle was rejected: the Java SDK no longer
ships a tests jar and the publishing plugin is newer). The NuGet package now carries a readme, copyright and release notes.
Use the exact version (a release candidate is not picked up by a plain install); report problems to support@openroam.com.

## 0.1.0-rc.1

First release candidate of the OpenRoam Partner API SDKs for TypeScript/JavaScript, Python, Go, PHP, Java, .NET, Dart and Ruby,
generated from the Partner API description. A release candidate is not picked up by a plain install: use the exact version (see each
README). The stable 0.1.0 follows once the candidate has been tried; report problems to support@openroam.com.
