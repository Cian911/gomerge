### go-merge

#### Description
Gomerge is a tool to quickly merge several pull requests from your terminal. The intention of this tool is to simplfy, and eventually automate the merging of github pull requests. This tool should be able to run on most systems.

#### Requirements

You must have created a github personal access token (PAT) to use this tool. For information on how to do so, you can follow the documentation https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token

#### Install

To install the compiled binary, you can go to the releases tab, and download the version most suitable to your environment. Otherwise, you can download the latest binary like so:

```bash
wget wget https://github.com/Cian911/gomerge/releases/download/1.0.0/gomerge_1.0.0_linux_amd64.tar.gz
sudo tar -xvf gomerge_1.0.0_linux_amd64.tar.gz -C /usr/local/bin/
sudo chmod +x /usr/local/bin/gomerge
```
#### Usage

Below denotes the available commands and flags on the `gomerge` tool.

```bash
gomerge makes it simple to merge an open pull request from your terminal

Usage:
  gomerge [command]

Available Commands:
  help        Help about any command
  list        List all open pull request for a repo.

Flags:
  -h, --help           help for gomerge
  -r, --repo string    Pass name of repository as argument (organization/repo).
  -t, --token string   Pass your github personal access token (PAT).

Use "gomerge [command] --help" for more information about a command.
```

To get a list of open and active pull requests for a given repo, you can run the following command.

**N.B**: Please ensure to add your organization followed by the name of your repository. In most cases this will be your github username, but if referencing a repository that exists within an organization you have access to, be sure to substitute it for that E.G `google/example-repo`.

```bash
gomerge list -r Cian911/go-merge  -t ${GITHUB_TOKEN}
```

If there are any active and open pull requests for your given repository, you will see an output similar to below.

![gomerge Sample Output](https://i.imgur.com/UIsiEGd.png)

From here, follow the instructions to select which pull request you wish to merge, and hit enter. Your pull request should now have been merged, and you should get a similar message to below.

```bash
PR #3: Pull Request successfully merged.
```
