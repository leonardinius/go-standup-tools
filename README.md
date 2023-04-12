# go-jira-tools

go-jira-tools is a command-line tool written in Go that generates daily standup notes based on Jira activity feed. 
It extracts the relevant information from the activity feed and formats it into an easy-to-read template. 
Additionally, it automatically copies the output to the clipboard.

## Installation

To install go-jira-tools, follow the steps below:

1. Clone this repository to your local machine.
2. Copy the `run.example.sh` file to `run.sh` and adjust it as necessary.
3. Run `./run.sh`.

## Usage

To generate a daily standup notes template, run the following command:

```bash
./run.sh

```

This will produce the standup notes in the following format:

```
Copy code
2023/04/12 17:46:38 [INFO ] Text:

-   **2023-03-30**
-   retest OSA-3906 - Testing demo
-   closed OSA-3905 - test demo
-   OSA-3935 - Upgradeability in standby
-   OSA-3937 - Playing MP3
-   OSA-3936 - SWT TEST
-   OSA-3938 - Test
-   OSA-3939 - Negativ-Test: Feldlängenüberschreitung
-   OSA-1381 - Provider -Payment -Check scan-off
-   OSA-1635 - Auszahlung nach 4-Augenprinzip
    2023/04/12 17:46:38 [INFO ] ✔️ copied to clipboard!
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are always welcome! Please feel free to open an issue or submit a pull request.

## Contact

If you have any questions or concerns, please feel free to contact me via Github issue.
