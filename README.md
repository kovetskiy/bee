# bee

A simple program written on a knee that highlights last line if there is no
input after a second.

# Example
```
bash -c 'while true; do echo 1; sleep 0.05; echo 2; sleep 1.1; echo 3; sleep 0.01; done' | ./bee
```

Bee will make line `2` yellow after a second because there is a pause more than
1 second between `2` and `3`.
