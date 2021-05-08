## gomerge
![GitHub Actions Status](https://github.com/Cian911/gomerge/workflows/Release/badge.svg)

<p align="center">
  <img style="float: right;" src="sample/gomerge-gopher.png" alt="Gomerge logo"/>
</p>

#### Description
Gomerge is a tool to quickly enable you to bulk merge Github pull requests from your terminal. The intention of this tool is to simplfy, and eventually automate the merging of github pull requests. This tool should be able to run on most systems.

#### Requirements

You must have created a github personal access token (PAT) to use this tool. For information on how to do so, you can follow the documentation https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token

![Gomerge Demo](https://i.imgur.com/2vX6ty3.gif)

#### Install 

To install the latest version via `homebrew`, please run the following.

```bash
brew tap Cian911/gomerge
brew install gomerge

// Check the binary is working as expected.
gomerge -h
```

To install the compiled binary, you can go to the releases tab, and download the version most suitable to your environment. An example of such is below.

```bash
wget https://github.com/Cian911/gomerge/releases/download/1.2.0/gomerge_1.2.0_Linux_x86_64.tar.gz
sudo tar -xvf gomerge_1.2.0_Linux_x86_64.tar.gz -C /usr/local/bin/
sudo chmod +x /usr/local/bin/gomerge
```

###### Upgrade

If you have the tool installed already via homebrew, you can upgrade by running the following:
```bash
brew upgrade gomerge
```

#### Usage

Below denotes the available commands and flags on the `gomerge` tool.

```bash
Gomerge makes it simple to merge an open pull request from your terminal.

Usage:
  gomerge [command]

Available Commands:
  help        Help about any command
  list        List all open pull request for a repository you wish to merge.

Flags:
  -c, --config string   Pass an optional config file as an argument with list of repositories.
  -h, --help            help for gomerge
  -r, --repo string     Pass name of repository as argument (organization/repo).
  -t, --token string    Pass your github personal access token (PAT).

Use "gomerge [command] --help" for more information about a command.
```

To get a list of open and active pull requests for a given repo, you can run the following command.

**N.B**: Please ensure to add your organization followed by the name of your repository. In most cases this will be your github username, but if referencing a repository that exists within an organization you have access to, be sure to substitute it for that E.G `google/example-repo`.

```bash
gomerge list -r Cian911/go-merge -t ${GITHUB_TOKEN}
```

If there are any active and open pull requests for your given repository, you will see an output similar to below.

![gomerge Sample Output](https://i.imgur.com/UIsiEGd.png)

From here, follow the instructions to select which pull request you wish to merge, and hit enter. Your pull request should now have been merged, and you should get a similar message to below.

```bash
PR #3: Pull Request successfully merged.
```

##### Bulk Merging Pull Requests

As of version `1.1.0` there is a new option available to pass a config.yaml as an arugment to the `gomerge` tool which will give the user the option to configure a list of repositories in order to more easily _bulk merge_ pull requests.

You should first create a `config.yaml` file in the following format.

```yaml
organization: Cian911
repositories:
- pr-test
- syncwave
```

You can then run the tool like so, passing the config file as a flag.

```bash
gomerge list -t $GITHUB_TOKEN -c config.yaml
```

You should see a list of active and open pull requets from the repositories you have defined in your configuration file.

![gomerge Bulk Output Sample](https://imgur.com/zROhCYV.png)
