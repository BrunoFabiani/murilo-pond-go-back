import http from 'k6/http'
import { check, sleep } from 'k6'

export const options = {
    stages: [
      { duration: '30s', target: 10 },
      { duration: '30s', target: 50 },
      { duration: '30s', target: 100 },
      { duration: '30s', target: 0 },
    ],
    
};


export default function () {

    const url = 'http://localhost:8080/telemetry';

    const payload = JSON.stringify({
        device_id: `device-${__VU}`, //$ insert 
        timestamp: new Date().toISOString(), //js data method converts time to standadrd UTC string
        sensor_type: 'temperature',
        reading_nature: 'analog',
        value: 23.5,
    });
    const params = { headers: { 'Content-Type': 'application/json' } };

    const res = http.post(url, payload, params);
    check(res, { 'status is 201': (r) => r.status === 201 });
    
  }
