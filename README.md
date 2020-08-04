# NE0FS

free storage on neo network inspired by the flood attack.

# Usage

## Upload

```sh
./NE0FS -c upload -f data.jpg
```

get the printed `<headerhash>` which can be used to download

## Download

```sh
./NE0FS -c download -s <headerhash> -f dst.jpg
```

## Detect Node

```sh
./NE0FS -c detectnode -f nodes.json
```

then you can use `./NE0FS -c download -s <headerhash> -f dst.jpg -nodelist nodes.json` to download faster.

# TODO

- multi-threading upload / download

- use p2p message instead of rpc