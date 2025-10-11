import http from "k6/http";
import { check } from "k6";

/*
    going on maximum performance to check how system degrade
*/

export let options = {
    stages: [
        { duration: "20s", target: 150 },
        { duration: "20s", target: 300 },
        { duration: "20s", target: 500 },
        { duration: "3m", target: 500 }
    ]
};

export default function() {
    let res = http.get("http://avito-shop-service:8080/api/info");
    check(res, { "status is 200": (r) => r.status === 200 });
}