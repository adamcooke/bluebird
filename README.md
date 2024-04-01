# Bluebird 🐦

Bluebird is a TUI application to allow you to search through a number of hierarchical commands and execute them quickly. Ideal for executing SSH sessions to servers or running complex commands on Kubernetes.

![Screenshot](https://github.com/adamcooke/bluebird/assets/4765/b59ac886-995d-4d0e-a58e-1cdd296e5bed)

## Usage

To get started, you'll need a config file. By default, this lives in `~/.bluebird.yaml`. This file contains all the commands that you want to be able to access through the `bluebird` tool.

### Example configuration

The first (and only) thing you need to configure is what items you want to be accessible. This is done by providing an array of items. At present, an "item" can be one of two types: a list or a command. A list is sub-list which allows you to group commands and a command is something that will be executed when it is selected.

```yaml
items:
  - name: My Servers
    description: SSH to my servers
    emoji: ⌨️
    type: list
    listOptions:
      items:
        - name: Primary server
          type: command
          commandOptions:
            command: ssh primary.example.com
        - name: Secondary server
          type: command
          commandOptions:
            command: ssh secondary.example.com
```
