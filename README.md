# AWS Organizations Visualiser

![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/CentricaDevOps/aws-organizations-visualiser/unit-test.yml?label=Unit%20Tests)
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/CentricaDevOps/aws-organizations-visualiser/golangci-lint.yml?label=Golang%20Linting)
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/CentricaDevOps/aws-organizations-visualiser/update-release.yml?label=Release)


This is a tool that can be used to visualise the structure of an AWS
Organizations structure. It can be used to generate a JSON representation of the
structure or to display the structure in the CLI.

## Usage

This tool generates a structure that represents the AWS Organizations structure
and then, based on the flags passed in, either displays the structure in the CLI
or outputs the structure to a JSON file or both.

To use this tool, download the latest release for your platform from the
[releases page](https://github.com/CentricaDevOps/aws-organizations-visualiser/releases)
and then run the binary with the required flags. Please note that for this tool
to work you must be logged into an AWS account that has access to the AWS
Organizations API and you must have the correct permissions to access the
required information.

If you wish to build the tool from source, clone the repository and run the
following command:

    make build

This will build the tool for your current platform and place the binary in the
`bin` directory.

If you wish to simply run the tool without building it, run the following
command:

    make run

This will run the tool for your current platform.

Finally, if you wish to lint and test the tool, run the following command:

    make test

This will run the linter and unit tests for the tool.

### Example

To generate a JSON representation of the AWS Organizations structure and output
it to a file called `output.json` run the following command:

    aws-organizations-visualiser -include-visual=false

To generate a JSON representation of the AWS Organizations structure and output
it to a file called `output.json` and display the structure in the CLI run the
following command:

    aws-organizations-visualiser -o output.json

To generate a visual representation of the AWS Organizations structure in the 
CLI run the following command:

    aws-organizations-visualiser -include-json=false


### Flags

Usage:

    aws-organizations-visualiser [flags]

Flags:

    -remove-suspended-accounts
		Remove suspended accounts from the output (default false)
    -include-json
        Include the JSON representation of the AWS Organizations structure in the output (default true)
    -include-visual
        Include the visual representation of the AWS Organizations structure in the output (default true)
    -o string
        The output file for the JSON representation of the AWS Organizations structure (default "output.json")

## Contributing

If you have any suggestions or issues, please raise them in the issues section
of this repository. If you wish to contribute to this project, please fork the
repository and raise a pull request.

## License

This project is licensed under the GNU GPLv3 License - see the
[LICENSE](LICENSE) file for details.