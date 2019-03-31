import axios from "axios";
import { SerPrefix } from "../config";
const qs = require('qs');
const http = axios.create({
    method: 'post',
    baseURL: `http://${ SerPrefix }`,
    headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
    },
    timeout: 15000,
});

export const request = async ({
    url,
    data,
    method = 'POST'
}) => {
    method = method.toLocaleUpperCase();
    if (method === 'GET') {
        url += `?${qs.stringify(data)}`;
    } else {
        data = qs.stringify(data);
    }
    return http({
        url: url,
        method,
        data,
    })
};