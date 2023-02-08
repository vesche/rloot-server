## rloot file system structure

```
//
// fs structure
//
// opt/
//     rloot/
//         _meta/
//             rloot.db
//             images/
//                 <poster files>
//         movies/
//             <media files>
//         tv/
//             <media files>
//
```

## Server commands

- `rloot server --install`
    - cross-platform service installer
    - Ask some options, make a config file

- `rloot server --start --config $HOME/.rloot`
    - start the server
    - default location for config?

- `rloot server --stop`

- `rloot server --uninstall`

- `rloot metadata`
    - Run/refresh the server metadata files/db
    - Should also be able to hit this from an API endpoint

- `rloot add movie <file_path>`
    - Add new file, grab meta-data

- `rloot remove ?`

## Server start/stop cross-platform

- Linux
    - systemd / `systemctl start rloot`
- macOS
    - LaunchAgent / `launchctl`
- Windows
    - Service / `net start rloot`

## Notes

- I'm not interested in supporting multiple libraries of the same type
