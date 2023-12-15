# Bear OJ(Online Judge) Backend

Bear OJ(Online Judge) is for private competive programming contests, based on open-source online-judge system, [DMOJ](https://github.com/DMOJ/online-judge). Though Bear OJ is based on DMOJ, the system besides judging system is fully customized.

## Features

Bear OJ provides monitoring system, which prohibits those actions by actively monitering participants.

- Prohibits searching, viewing another resources online or offline by
    - Only allowing full-screen participation
    - Alerting the contest administrator when the focus has been lost from the website

## Technologies

Bear OJ Backend is written in [Go](https://go.dev/), using [Gin Gonic](https://github.com/gin-gonic/gin) HTTP library.
