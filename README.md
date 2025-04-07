# **GoShell**

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](link) [![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

## Introduction

```
  ,----..              .--.--.     ,---,                ,--,    ,--,   
 /   /   \            /  /    '. ,--.' |              ,--.'|  ,--.'|   
|   :     :    ,---. |  :  /`. / |  |  :              |  | :  |  | :   
.   |  ;. /   '   ,'\;  |  |--`  :  :  :              :  : '  :  : '   
.   ; /--`   /   /   |  :  ;_    :  |  |,--.   ,---.  |  ' |  |  ' |   
;   | ;  __ .   ; ,. :\  \    `. |  :  '   |  /     \ '  | |  '  | |   
|   : |.' .''   | |: : `----.   \|  |   /' : /    /  ||  | :  |  | :   
.   | '_.' :'   | .; : __ \  \  |'  :  | | |.    ' / |'  : |__'  : |__ 
'   ; : \  ||   :    |/  /`--'  /|  |  ' | :'   ;   /||  | '.'|  | '.'|
'   | '/  .' \   \  /'--'.     / |  :  :_:,''   |  / |;  :    ;  :    ;
|   :    /    `----'   `--'---'  |  | ,'    |   :    ||  ,   /|  ,   / 
 \   \ .'                        `--''       \   \  /  ---`-'  ---`-'  
  `---`                                       `----'                   
```

Welcome to **GoShell**, an ambitious and dynamic shell application written in Go. This project aims to bring robust, modular, and user-friendly shell commands to your fingertips. Whether you’re managing files or navigating directories, GoShell provides an intuitive command-line interface to streamline your workflows.

## Features

- **Concatenate and Display File Contents**: Use `cat` to view and merge files with additional options.
- **Change Directory**: The `cd` command allows seamless navigation between directories.
- **Remove Files and Directories**: Efficiently remove files or directories with various flags.
- **Modular Architecture**: Built with a clean structure for easy maintenance and expansion.

## Installation

To get started with GoShell, follow these steps:

1. **Clone the Repository:**

    ```bash
    git clone https://github.com/IshaanNene/GoShell.git
    ```

2. **Navigate to the Project Directory:**

    ```bash
    cd GoShell
    ```

3. **Build the Application:**

    ```bash
    go build -o goshell ./cmd/goshell
    ```

## Usage

Here are some examples of how to use GoShell:

- **Change Directory:**

    ```bash
    ./goshell cd /path/to/directory
    ```

- **Concatenate Files:**

    ```bash
    ./goshell cat n file1.txt file2.txt
    ```

- **Remove a File:**

    ```bash
    ./goshell rm file.txt
    ```

## Contributing

Contributions are welcome! Please follow these guidelines to contribute:

1. **Fork the Repository**: Create a personal copy of the repository on GitHub.
2. **Create a Branch**: Make your changes in a separate branch.
3. **Submit a Pull Request**: Describe your changes and submit a pull request for review.

For more details, check out our [CONTRIBUTING.md](CONTRIBUTING.md).

## License

GoShell is licensed under the [MIT License](LICENSE). See the LICENSE file for more information.

## Acknowledgments

- Thanks to the Go community for the ongoing support and contributions.
- Special thanks to the maintainers of the `cobra` package for their excellent work on command-line interfaces.

## Final Thoughts

I hope GoShell makes your command-line experience more enjoyable and efficient. Try it out, and let us know what you think!

```
            ___
        .-"; ! ;"-.
      .'!  : | :  !`.
     /\  ! : ! : !  /\
    /\ |  ! :|: !  | /\
   (  \ \ ; :!: ; / /  )
  ( `. \ | !:|:! | / .' )
  (`. \ \ \!:|:!/ / / .')
   \ `.`.\ |!|! |/,'.' /
    `._`.\\\!!!// .'_.'
       `.`.\\|//.'.'
        |`._`n'_.'| 
        "----^----"
```

