# Install and Setup
- `Go 1.16x`version install https://golang.org/dl/
- `GOPATH` and `GOROOT` setup after Go installation
    - I'm using default `GOROOT` path
    - GOPATH can be any place you want to place go lang projects. For example, I created `go_repos` directory in `~/Documents` and set my `GOPATH` as `~/Documents/go_repos` in my `env`
    - Have a `src` directory in `GOPATH`. 
    ```
    # in my local file structure
    - ~/ 
      - Documents/
        - go_repos/
          - src/
    ```
- This project isn't using `go module` so set gomudle off by `$export GO111MODULE = off` and create `personal/` directory in `src/` in order to prevent import path error in the future for this proejct
- clone this project in the `personal`
```
    - ~/ 
      - Documents/
        - go_repos/
          - src/
            - personal
              - valet-parking-lot
```

# Build
```
in valet-parking-lot/
$go build valetparking.go #creates compiled executable file valetparking

#run program
$./valetparking <your_txt_input_file> # e.g) ./valetparking input1.txt
```
* Place your input file in `valet-parking-lot/inputs` directory


# Dcoumentation
https://docs.google.com/document/d/1UIENDkIf9Mt0e4PXxngq5VT6_UezTejJ5hImtyHehsQ/edit?usp=sharing
