## gomerge
![GitHub Actions Status](https://github.com/Cian911/gomerge/workflows/Release/badge.svg) ![GitHub Actions Status](https://github.com/Cian911/gomerge/workflows/Test/badge.svg) ![Downloads](https://img.shields.io/github/downloads/cian911/gomerge/total.svg)

<p align="center">
  <img style="float: right;" src="sample/gomerge-gopher.png" alt="Gomerge logo"/>
</p>

#### Description
Gomerge is a tool to quickly enable you to bulk merge and/or approve Github pull requests from your terminal. The intention of this tool is to simplfy merging large numbers of Github pull requests, as they can certainly grow when you're maintaining a large number of repositories (Ahem, dependabot!). This tool should be able to run on most systems.

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
wget https://github.com/Cian911/gomerge/releases/download/3.4.0/gomerge_3.4.0_Linux_x86_64.tar.gz
sudo tar -xvf gomerge_3.4.0_Linux_x86_64.tar.gz -C /usr/local/bin/
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
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        List all open pull request for a repository you wish to merge.
  version     Prints the current version and build information.

Flags:
  -a, --approve                      Pass an optional approve flag as an argument which will only approve and not merge selected repos.
      --close                        Pass an optional argument to close a pull request.
  -c, --config string                Pass an optional config file as an argument with list of repositories.
  -d, --delay int                    Set the value of delay, which will determine how long to wait between mergeing pull requests. Default is (6) seconds. (default 6)
  -e, --enterprise-base-url string   For Github Enterprise users, you can pass your enterprise base. Format: http(s)://[hostname]/
  -h, --help                         help for gomerge
  -l, --label stringArray            Pass an optional list of labels to filter pull requests. (label1,label2,label3)
  -m, --merge-method string          Pass an optional merge method for the pull request (merge [default], squash, rebase).
  -r, --repo string                  Pass name of repository as argument (organization/repo).
  -s, --skip                         Pass an optional flag to skip a pull request and continue if one or more are not mergable.
  -t, --token string                 Pass your github personal access token (PAT).

Use "gomerge [command] --help" for more information about a command.

```

To get a list of open and active pull requests for a given repo, you can run the following command. Note, this will list all available PRs for merging, if you want to just approve a list of PRs, then add the `-a` to the below command too.

**N.B**: Please ensure to add your organization followed by the name of your repository. In most cases this will be your github username, but if referencing a repository that exists within an organization you have access to, be sure to substitute it for that E.G `google/example-repo`.

```bash
gomerge list -r Cian911/gomerge -t ${GITHUB_TOKEN}
```

Something to note, if you have your github token defined as an environment variable in the following format `GITHUB_TOKEN`, gomerge will try and use that automagically if you choose not to define a token via a flag. Thanks to [@caioeverest](https://github.com/caioeverest) for this work!

If there are any active and open pull requests for your given repository, you will see an output similar to below.

![gomerge Sample Output](https://i.imgur.com/UIsiEGd.png)

From here, follow the instructions to select which pull request you wish to merge, and hit enter. Your pull request should now have been merged, and you should get a similar message to below.

```bash
PR #3: Pull Request successfully merged.
```

##### Bulk Merging/Approving/Closing Pull Requests

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

And again, if you just wish to approve a list of available PRs, just add the `-a` flag like so.

```bash
gomerge list -t $GITHUB_TOKEN -c config.yaml -a
```

And finally, if you wish to close a list of PRs, just add `--close`

```bash
gomerge list -t $GITHUB_TOKEN -c config.yaml --close
```

You should see a list of active and open pull requets from the repositories you have defined in your configuration file.

![gomerge Bulk Output Sample](https://imgur.com/zROhCYV.png)
