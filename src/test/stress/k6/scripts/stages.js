import http from "k6/http";
import { check } from "k6";


/*
 * Stages (aka ramping) is how you, in code, specify the ramping of VUs.
 * That is, how many VUs should be active and generating traffic against
 * the target system at any specific point in time for the duration of
 * the test.
 */ 

/*
    Increasing and decreasing VUs
*/

export let options = {
    stages: [
        { duration: "20s", target: 150 },
        { duration: "20s", target: 300 },
        { duration: "20s", target: 700 },
        { duration: "20s", target: 1000 },
        { duration: "20s", target: 1300 },
        { duration: "20s", target: 1600 },
        { duration: "20s", target: 2000 },
        { duration: "40s", target: 2000 },
        { duration: "20s", target: 1600 },
        { duration: "20s", target: 1300 },
        { duration: "20s", target: 1000 },
        { duration: "20s", target: 700 },
        { duration: "20s", target: 300 },
        { duration: "20s", target: 0 },
    ]
};

export default function() {
    let res = http.get("http://auth:8080/ping");
    check(res, { "status is 200": (r) => r.status === 200 });
}