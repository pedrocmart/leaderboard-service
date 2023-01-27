# Leaderboard Service

## Table of Contents
- [Description of the task](#description)
- [Solution](#solution)
    - [Requirements](#requirements-to-install)
    - [Build local development](#buildlocal)
    - [Running Tests](#tests)
- [APIs](#APIs)
    - [POST user/{user_id}/score](#post)
        - [Absolute](#postabsolute)
        - [Relative](#postrelative)
    - [GET ranking?type={type}](#get)
        - [Absolute](#getabsolute)
        - [Relative](#getrelative)
- [Environment Variables](#environment)

-----------------------






# Description
Your task is to implement a service that maintains a "real-time" ranking of users playing a specific game. 

Requirements:

<ul>
	<li>Clients submit user scores when they achieve important milestones in the game&nbsp;</li>
	<li>Clients can submit absolute scores or relative scores, for example: {"user" :123, "total": 250}, or {"user": 456, "score": "+10"}, or {"user": 789, score: "-20"}</li>
	<li>Any client can request the ranking at any time, using one of the following requests:</li>
	<li>Absolute ranking, for example: Top100, Top200, Top500</li>
	<li>Relative ranking, for instance: At100/3, meaning 3 users around position 100th of the ranking, that is positioned 97th, 98th, 99th, 100th, 101st, 102nd, 103rd</li>
</ul>

<p>We propose that it have the following endpoints:</p>
<ul>
	<li>[POST] user/{user_id}/score</li>
	<li>[GET] ranking?type=top100</li>
	<li>[GET] ranking?type=At100/3</li>
</ul>

<h3>Technical requirements</h3>
<ul>
	<li>We prefer&nbsp;you to develop the test in one of the main&nbsp;languages we use at SP (PHP or Golang)&nbsp;but if you don't&nbsp;feel comfortable enough with them, feel free to choose any other</li>
	<li>
<strong>Do&nbsp;not couple yourself to a&nbsp;specific framework</strong>, as the test is pretty simple we prefer to see how big is your knowledge of what goes under the hood</li>
	<li>
<strong>Do not use any database or external storage system</strong>, just keep the ranking <strong>in-memory</strong> (NB if you use a stateless language&nbsp;there's no need to keep this storage anywhere after the process dies)</li>
	<li>
<strong>The code must work and sort</strong> as specified in this document</li>
</ul>

<h3>Goals to evaluate</h3>

<ul>
	<li>How you approach the project (we left some stuff &nbsp;intentionally open, so you have to evaluate trade-offs and make some decisions)</li>
	<li>How you design and architecture the system&nbsp;</li>
	<li>How you test and ensure the overall quality of the solution</li>
</ul>


<h3>Nice to have</h3>
<ul>
	<li>Docker with docker-compose</li>
	<li>Readme to explain how we can check the code</li>
	<li>Testing with full code coverage or important code for you</li>
</ul>

----------


# Solution

<a id="requirements-to-install"></a>
## Requirements to install
  * Go 1.18
  * Docker
  * docker-compose

<a id="buildlocal"></a>
## Build for `Local` development
`make local`

<a id="tests"></a>
## Running Tests
 `make test`
______________
<a id="APIs"></a>
## APIs
<a id="post"></a>
### **[POST] user/{user_id}/score**

Since we can have two different submissions for the user (absolute or relative), the body must contain only one option, as described below:
<a id="postabsolute"></a>
#### **Absolute Score:**
Must be sent with the integer `total`. This will set the user's total score.

If the `user_id` sent didn't exist in the database, a new user will be created.

**Example:**

`[POST]` http://0.0.0.0:8894/user/1/score

`[JSON Body]`
```
{
    "total": 100
}
```

<a id="postrelative"></a>
#### **Relative Score:**
Must be sent with the string `score`. This will add or subtracts from user's score, depending on the first character sent. The API only accepts `score` which starts with `+` or `-`.

The user can have a negative score.

If the `user_id` sent didn't exist in the database, a new user will be created.

**Example:**

`[POST]` http://0.0.0.0:8894/user/1/score

`[JSON Body]`
```
{
    "score": "+20"
}
```

Response:

For both scenarios, the response will be the same, in case of success. A JSON containing the `user_id`, and the current `score`.
```
{
    "user_id": 1,
    "score": 100
}
```





<a id="get"></a>
### **[GET] ranking?type={type}**
In order to request the ranking, you need to set what kind of ranking do you want to see: absolute or relative. That said, the API will only accept the followin types as a parameter:

<a id="getabsolute"></a>
#### **Absolute Ranking:**

In this scenario, you must use the type `top`, followed by the number of users that you want in the ranking. The number must be greater than 0.


**Example:**

Below we define `type` as `top3`, meaning that we want the 3 top players from our ranking.

`[GET]` http://0.0.0.0:8894/ranking?type=top3

Response:
```
{
    "ranking": [
        {
            "position": 1,
            "user_id": 3,
            "score": 452
        },
        {
            "position": 2,
            "user_id": 2,
            "score": 101
        },
        {
            "position": 3,
            "user_id": 1,
            "score": 5
        }
    ]
}
```


<a id="getrelative"></a>
#### **Relative Ranking:**

In this scenario, you must use the type `at`, followed by the position you want in the ranking, a `/`, and a number of users around that position. Both numbers must be greater than 0.

**Example:**

Below we define `type` as `at10/2`, meaning 2 users around position 10th of the ranking, that is positioned 8th, 9th, 10th, 11th, 12th.

`[GET]` http://0.0.0.0:8894/ranking?type=at10/2

Response:
```
{
    "ranking": [
        {
            "position": 7,
            "user_id": 12,
            "score": 34
        },
        {
            "position": 8,
            "user_id": 7,
            "score": 27
        },
        {
            "position": 9,
            "user_id": 9,
            "score": 24
        },
        {
            "position": 10,
            "user_id": 8,
            "score": 19
        },
        {
            "position": 11,
            "user_id": 4,
            "score": 2
        },
        {
            "position": 12,
            "user_id": 5,
            "score": -72
        }
    ]
}
```
Edge case:
If you set your type as `at1/3`, you'll see the the top 1 user, and the 3 users next to her/him.
_____________

<a id="environment"></a>
## Environment Variables

| ENV VAR           | DESCRIPTION                                           | DEFAULT                              |
| ----------------- | ----------------------------------------------------- | ------------------------------------ |
| HOST              | service host                                          | 0.0.0.0                              |
| PORT              | service port                                          | 8884                                 

---
