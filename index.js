const express = require('express')
const nacl = require('tweetnacl');
const parse = require('co-body')
const publicKey = 'dc2a2ef24d22c445bd5a81bab30219e7f1ebbaa8035513457cac4b145b32cdc3';
const app = express()

app.post("/", async function (req, res) {
    console.log('function start')
    const signature = req.get('X-Signature-Ed25519');
    const timestamp = req.get('X-Signature-Timestamp');
    const body = await parse.json(req, {returnRawBody: true});
    console.log('extracting headers and parsing body');

    const isVerified = nacl.sign.detached.verify(
        Buffer.from(timestamp + body.raw),
        Buffer.from(signature, 'hex'),
        Buffer.from(publicKey, 'hex')
    );

    console.log('determining if token is valid');

    if (!isVerified) {
        console.log('token is not valid');
        return res.status(401).end('invalid request signature');
    }

    console.log('token is valid');

    if (body.parsed['type'] === 1) {
        console.log('respond to discord ping')
        res.setHeader('Content-Type', 'application/json');
        res.end(JSON.stringify({'type': 1}));
    }

    res.status(200).end("end of function")
});

exports.bot = app