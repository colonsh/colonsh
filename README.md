## :: colonsh

> **A simple command-line helper that gives you colon-based aliases (:pd, :po, :pa) and lets you define project directories, repo-specific actions, and simple shortcuts — all from one JSON config file.**

## Installation
### Homebrew (Recommended)
This is the easiest and recommended way to install `colonsh` on macOS or Linux.
```bash
brew tap stephenbaidu/tap
brew install colonsh
```
Then initialize and reload your shell:
```bash
colonsh setup
source ~/.zshrc   # Reload your shell profile, eg:
```

### Alternative Installation (All Platforms)
Download a binary from the [GitHub Releases page](https://github.com/stephenbaidu/colonsh/releases) page and move it to a directory in your PATH:
- macOS/Linux → /usr/local/bin
- Windows → a custom bin folder

## Configuration
`colonsh` uses a single configuration file: **`~/colonsh.json`**. Open it with:
```bash
colonsh config
```
Example `colonsh.json`
```json
{
  "open_cmd": "code .",
  "aliases": [
    { 
      "name": "config", 
      "cmd": "code ~/colonsh.json" 
    },
    { 
      "name": "source", 
      "cmd": "source ~/.zshrc" 
    }
  ],
  "project_dirs": [
    { 
      "path": "~/Code/Work", 
      "exclude": ["archived"] 
    },
    { 
      "path": "~/Code/Personal" 
    }
  ],
  "git_repos": [
    {
      "slug": "stephenbaidu/colonsh",
      "actions": [
        { 
          "name": "Deploy locally", 
          "cmd": "go build . && mv colonsh $HOME/bin/colonsh" 
        },
        { 
          "name": "Run tests", 
          "cmd": "go test ./..." 
        }
      ]
    }
  ]
}
```

### Configuration Sections

### `open_cmd`
Defines the default command used by `:po` (Project Open) to open the current project. If not set, the default is: `code .`

### `aliases`

The **`aliases`** array defines simple custom commands accessible from anywhere in your shell via the `:` prefix (e.g., `:config`, `:source`). These are simple command substitutions that run shell commands.

| Key | Description |
| :--- | :--- |
| **`name`** | The specific alias name to be used after the colon (e.g., `:config`). |
| **`cmd`** | The raw shell command that `colonsh` executes when the alias is called. |

### `project_dirs`

The **`project_dirs`** array instructs `colonsh` where to scan for Git repositories on your system. This data is used by the `:pd` command to provide a searchable, quick-jump list of all your projects.

| Key | Description |
| :--- | :--- |
| **`path`** | The root directory path where `colonsh` should recursively look for Git repositories. Tilde (`~`) expansion is supported. |
| **`exclude`** | *(Optional)* A list of subdirectory names to ignore during the scan (e.g., excluding an `archived` folder within a large work directory). |

### `git_repos`

The **`git_repos`** array defines specific actions and behaviors for individual Git repositories. This is the most powerful section, enabling context-aware actions via the `:pa` command.

| Key | Description |
| :--- | :--- |
| **`slug`** | The unique identifier for the repository, typically in the format `organization/repo-name` (e.g., `stephenbaidu/colonsh`). |
| **`actions`** | A list of structured commands that only become available via `:pa` when your current working directory is inside this specific repository. |
| **`actions.name`** | The descriptive name displayed in the interactive list when running `:pa`. |
| **`actions.cmd`** | The shell command to be executed when this action is selected. |

***

## Development

### Building from Source

```bash
# Build the binary in the current directory
go build . 

# Run the compiled binary
./colonsh init zsh
```

### Testing

```bash
go test ./...
```

## License

`colonsh` is distributed under the **MIT License**. See the [LICENSE](LICENSE) file for details.