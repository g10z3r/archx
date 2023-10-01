#!/bin/bash
set -e

mongosh <<EOF
    use admin 

    db.createUser({
        user: "$MONGO_USER",
        pwd: "$MONGO_USER_PASSWORD",
        roles: [ 
            { 
                role: "dbAdminAnyDatabase", 
                db: "admin" 
            }, 
            "readWriteAnyDatabase"
        ]
    });
EOF