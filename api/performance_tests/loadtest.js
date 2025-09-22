import http from 'k6/http';
import { check, sleep } from 'k6';

let counter = 0; 

export let options = {
    vus: 50,          // Use 1 VU to control order
    duration: '30s',  
};

export default function () {
    // 1️ POST /api/shorten
    let payload = JSON.stringify({ long_url: 'https://example.com' });
    let headers = { 'Content-Type': 'application/json' };

    let postRes = http.post('http://host.docker.internal:9000/api/shorten', payload, { headers });

    // Check POST response
    check(postRes, {
        'POST status was 201': (r) => r.status === 201,
        'POST response time < 200ms': (r) => r.timings.duration < 200,
    });

    // Extract shortCode from response
    let shortURL = '';
    try {
        let body = JSON.parse(postRes.body);
        shortURL = body.short_url;
        check(body, {
            'POST response has short_url': (b) => b.short_url !== undefined,
        });
    } catch (e) {
        console.error('Failed to parse POST response:', e);
        return;
    }

    let shortCode = shortURL.split('/').pop();

    //  Check if shortCode is numeric and starts from "0"
    check(shortCode, {
        [`shortCode is ${counter}`]: (code) => code === String(counter),
    });

    // 2️ GET /:shortCode
    let getRes = http.get(`http://host.docker.internal:9000/${shortCode}`, { redirects: 0 });

    check(getRes, {
        'GET status was 302': (r) => r.status === 302,
        'GET response time < 10ms (cache hit)': (r) => r.timings.duration < 10,
    });

    counter++; 

    sleep(1);
}
