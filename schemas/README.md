# Butane-Schemas

Butane Schemas that helps creating Butane files for [Fedora Coreos](https://fedoraproject.org/fr/coreos/) and [Flatcar Linux](https://www.flatcar.org).

This json Schema that can be used as a helper to write a [butane config file](https://docs.fedoraproject.org/en-US/fedora-coreos/producing-ign/#_configuration_process) according to [its specifications](https://coreos.github.io/butane/specs/)

## Use it with VsCode

- Install [VS Code](https://github.com/microsoft/vscode)
- Install the [Red hat YAML extension](https://github.com/redhat-developer/vscode-yaml)
- Associate a schema in the YAML file `# yaml-language-server: $schema=<urlToTheSchema>` : [doc](https://github.com/redhat-developer/vscode-yaml?tab=readme-ov-file#associating-a-schema-in-the-yaml-file)
  - For butane schema, use `# yaml-language-server: $schema=https://coreos.github.io/butane/schemas/Butane-Schema.json"`

## Setup vs code to associate a schema to your *.bu files whitout setting schema manually

- Edit your `settings.json` file like :

New with version 1.8.0 of the [Red hat YAML extension](https://github.com/redhat-developer/vscode-yaml) and since this I've published this work on [JSON Schema Store](https://www.schemastore.org) you just need to change your files associations.

```JSON
"settings": {
  "files.associations": {
    "*.bu": "yaml"
  }
}
```

## Use it with Sublime text

You can also use these butane schemas with Sublime text editor thanks to [Joe Doss](https://github.com/jdoss) via [sublime-butane](https://github.com/jdoss/sublime-butane).

## List of clients that use uses the Red hat YAML extension

[some known clients consuming Red hat YAML extension](https://github.com/redhat-developer/yaml-language-server?tab=readme-ov-file#clients)
