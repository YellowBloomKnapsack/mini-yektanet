#!/usr/bin/env bash
scp -P 2233 panel/.env comp2@95.217.125.138:/home/comp2/mini-yektanet/panel/.env
scp -P 2233 adserver/.env comp2@95.217.125.138:/home/comp2/mini-yektanet/adserver/.env
scp -P 2233 common/.env comp2@95.217.125.138:/home/comp2/mini-yektanet/common/.env
scp -P 2233 eventserver/.env comp2@95.217.125.138:/home/comp2/mini-yektanet/eventserver/.env
scp -P 2233 publisherwebsite/.env comp2@95.217.125.138:/home/comp2/mini-yektanet/publisherwebsite/.env