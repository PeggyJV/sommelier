I couldn't find a way to spin up the chain, using their version number. I've initated an upgrade module, containing the `go.mod` and `go.sum` files. Here's how this suite works:

- Seclude suite to conform previous chain version.
- Run upgrade to current version.
- Interact with upgraded version to ensure upgrade was successful.

## Usage For Testing New Upgrades
- Change Sommelier and Gravity versions in `go.mod`.
- Work on suite setup files to reflect changes in Sommelier and Gravity.
- Build docker image of the upgrade version.