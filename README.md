# Installation

[Release v1.0.0 · Aivyss/password-manager · GitHub](https://github.com/Aivyss/password-manager/releases/tag/v1.0.0)

1. Download the program that matches your OS and architecture from the provided link.

2. Rename the downloaded program to `pwm`.

3. Move the downloaded file to an appropriate directory and set the environment variable:

```shell
export PATH={/your/directory}:$PATH
```

# Execution Test (Check Version)

If you have set the environment variable, you can run the password manager from any directory!

## Check the version

Command:

```shell
pwm version
```

```shell
pwm v
```

Result:

```shell
[pwm] version: v1.0.0
```

If you see the above result, the application is working correctly.

# User Registration

The password manager allows you to register each user and manage passwords for each user.
NOTE: This is separate from the operating system's users.

Command:

```shell
pwm user create -name {user name} -pw {user password}
```

This command does not return any message upon success. If no error message is returned, the user creation was successful.

# User Login (Console Usage)

Once a user is registered, they can log in.

Command:

```shell
pwm user login -name {user name} -pw {user password}
```

Result:

```shell
[pwm] you entered password-manager console 
[pwm][main console] >
```

If the login is successful, you will enter the console as shown above. From this console, you can register, retrieve, and update passwords.

# Registering a Password

Password registration is managed like a key-value store. Even if the program is closed, the passwords will not disappear, and you can check them again when you log in. Passwords are stored in a local database and encrypted, making the original values unreadable to anyone except the user.

Command:

```shell
[pwm][main console] > save -k {your key} -pw {your password}
```

If this command is successful, it will not return any result. If there is a problem, it will return a message.

# Viewing the List of Registered Passwords

Command:

```shell
[pwm][main console] > list
```

Result (simple example):

```shell
  KEY  |         CREATED DATE          |       LAST UPDATED DATE
-------+-------------------------------+--------------------------------
  abcd | 2023-11-18 14:50:14 +0000 UTC | 2023-11-18 14:50:14 +0000 UTC
```

# Checking a Password

To check a password, use the following command.

Command:

```shell
[pwm][main console] > get -k {your key}
```

Result:

```shell
[pwm][main console] check your password: {your password} 
[pwm][main console] please 'clear' command after checking your password
```



# Clearing the Console Display

Command:

```shell
[pwm][main console] > clear
```

# Updating a Password

To update a password, use the following command.

Command:

```shell
[pwm][main console] > update -k {your key} -pw {your new password}
```

Result:

```shell
[pwm][main console] please enter your master password again:
```



This command requests the user's ID again. If the user password matches, the password is changed. If the password change is successful, it does not return any result.

# Exiting the Console

Command:

```shell
[pwm][main console] > exit
```








