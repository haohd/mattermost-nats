# Mattermost NATS
## Build
### Install npm 10
```
npm install -g npm@10
```

### Clone docconv to `./deps/docconv`

```
git clone https://github.com/sajari/docconv.git
```

### Edit `./mattermost/server/go.mod`

- Add below line to the head of the file:

```
replace code.sajari.com/docconv/v2 => ../../deps/docconv
```

### Install `gpg` tool (must need on Mac)

```
brew install gpg
#brew install gnupg
```

### Build MM

```
./run.sh build_mm
```

### Build Docker

```
./run.sh build_docker
```