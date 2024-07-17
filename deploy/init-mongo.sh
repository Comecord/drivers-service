#!/bin/bash

mongosh admin --host localhost -u root -p demoglonass2024CRM <<EOF
use demoglonass
db.createUser({
  user: "admin",
  pwd: "demoglonass2024CRMadmin",
  roles: [{ role: "readWrite", db: "demoglonass" }]
})
EOF