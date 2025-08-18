# age-vanity-keygen

This tool generates [age](https://github.com/FiloSottile/age) X25519 identity with a recipient that has a specified prefix.
The output is identical to `age-keygen`.

Compared to [similar tools](#similar-tools), it uses the fastest search algorithm, see [vanity25519](https://github.com/AlexanderYastrebov/vanity25519) for algorithm implementation 🚀

## Usage

Install the tool locally and run:
```console
$ go install github.com/AlexanderYastrebov/age-vanity-keygen@latest

$ age-vanity-keygen 23456
Found age123456... in 0s after 15855390 attempts (43267922 attempts/s)
# created: 2025-08-18T18:18:18+02:00
# public key: age123456gpgacec4alqvqnfdacx6djhx98wzwn4l3eh5q5n5ec2evdsfzn7tn
AGE-SECRET-KEY-1XRTF5T02CR2HEC29RAH29Y46DPHQ7EAPK5EEPYKTFE3682LPWSCS4CXJSX

$ echo AGE-SECRET-KEY-1XRTF5T02CR2HEC29RAH29Y46DPHQ7EAPK5EEPYKTFE3682LPWSCS4CXJSX | age-keygen -y
age123456gpgacec4alqvqnfdacx6djhx98wzwn4l3eh5q5n5ec2evdsfzn7tn
```

or use the Docker image:
```console
$ docker pull ghcr.io/alexanderyastrebov/age-vanity-keygen:latest
$ docker run  ghcr.io/alexanderyastrebov/age-vanity-keygen:latest 23456
```

## Performance

The tool checks ~40'000'000 keys per second on a laptop.
In practice, it finds a 6-character prefix within a minute.
Each additional character increases search time by a factor of 32.

## Similar tools

* [vanity-rage](https://github.com/siltyy/vanity-rage)
* [vanity-age](https://github.com/yawning/vanity-age)
* [vanity-age-keygen](https://codeberg.org/RachaelAva1024/vanity-age-keygen)
