import math
import sys

from dateutil.parser import parse
from pymongo import MongoClient

ZERO_TIMESTAMP = parse('0001-01-01T00:00:00.0Z')
HEIGHT = 12_310_005
TIMESTAMP = parse("2023-08-18T12:10:36.572027592Z")

db = MongoClient()["sentinelhub-2"]

cursor, height = db["blocks"].find().sort([("height", -1)]).limit(1), 0
for item in cursor:
    height = item["height"]

print("Latest height", height)
if height != HEIGHT - 1:
    print("Exiting...")
    sys.exit(1)

db["subscription_quotas"].rename("subscription_allocations", dropTarget=True)

db["deposits"].update_many(
    {},
    {
        "$rename": {
            "address": "addr",
        },
    },
)

# -------------------------------------------------------------------------------------------------------------------- #

db["events"].update_many(
    {},
    {
        "$rename": {
            "acc_address": "acc_addr",
            "allocated": "granted_bytes",
            "consumed": "utilised_bytes",
            "node_address": "node_addr",
            "price": "gigabyte_prices",
            "prov_address": "prov_addr",
        }
    },
)
db["events"].update_many(
    {},
    {
        "$unset": {
            "free": 1,
        }
    },
)
db["events"].update_many({"status": "STATUS_ACTIVE"}, {"$set": {"status": "active"}})
db["events"].update_many(
    {"status": "STATUS_INACTIVE_PENDING"}, {"$set": {"status": "inactive_pending"}}
)
db["events"].update_many(
    {"status": "STATUS_INACTIVE"}, {"$set": {"status": "inactive"}}
)

# -------------------------------------------------------------------------------------------------------------------- #

db["nodes"].update_many(
    {},
    {
        "$rename": {
            "address": "addr",
            "price": "gigabyte_prices",
        },
    },
)
db["nodes"].update_many(
    {},
    {
        "$unset": {
            "provider": 1,
        }
    },
)
db["nodes"].update_many({"status": "STATUS_ACTIVE"}, {"$set": {"status": "active"}})
db["nodes"].update_many({"status": "STATUS_INACTIVE"}, {"$set": {"status": "inactive"}})

cursor = db["nodes"].find()
for item in cursor:
    db["nodes"].find_one_and_update(
        {
            "addr": item["addr"],
        },
        {
            "$set": {
                "hourly_prices": item["gigabyte_prices"],
            },
        },
    )

# -------------------------------------------------------------------------------------------------------------------- #

db["plans"].update_many(
    {},
    {
        "$rename": {
            "provider_address": "prov_addr",
            "price": "prices",
            "bytes": "gigabytes",
            "validity": "duration",
            "node_addresses": "node_addrs",
            "add_height": "create_height",
            "add_timestamp": "create_timestamp",
            "add_tx_hash": "create_tx_hash",
        },
    },
)
db["plans"].update_many({"status": "STATUS_ACTIVE"}, {"$set": {"status": "active"}})
db["plans"].update_many({"status": "STATUS_INACTIVE"}, {"$set": {"status": "inactive"}})

cursor = db["plans"].find()
for item in cursor:
    db["plans"].find_one_and_update(
        {
            "id": item["id"],
        },
        {
            "$set": {
                "gigabytes": int(math.ceil(int(item["gigabytes"]) / 1e9)),
            },
        },
    )

# -------------------------------------------------------------------------------------------------------------------- #

db["providers"].update_many(
    {},
    {
        "$rename": {
            "address": "addr",
        },
    },
)

cursor = db["providers"].find()
for item in cursor:
    db["providers"].find_one_and_update(
        {
            "addr": item["addr"],
        },
        {
            "$set": {
                "status": "active",
                "status_at": HEIGHT,
                "status_timestamp": TIMESTAMP,
                "status_tx_hash": "",
            },
        },
    )

# -------------------------------------------------------------------------------------------------------------------- #

db["sessions"].update_many(
    {},
    {
        "$rename": {
            "subscription": "subscription_id",
            "address": "acc_addr",
            "node": "node_addr",
        },
    },
)
db["sessions"].update_many({"status": "STATUS_ACTIVE"}, {"$set": {"status": "active"}})
db["sessions"].update_many(
    {"status": "STATUS_INACTIVE_PENDING"}, {"$set": {"status": "inactive_pending"}}
)
db["sessions"].update_many(
    {"status": "STATUS_INACTIVE"}, {"$set": {"status": "inactive"}}
)

cursor = db["sessions"].find({"status": {"$ne": "inactive"}})
for item in cursor:
    db["sessions"].find_one_and_update(
        {
            "id": item["id"],
        },
        {
            "$set": {
                "end_height": HEIGHT,
                "end_timestamp": TIMESTAMP,
                "status": "inactive",
                "status_at": HEIGHT,
                "status_timestamp": TIMESTAMP,
            },
        },
    )
    db["events"].insert_one(
        {
            "type": "Session.UpdateStatus",
            "height": HEIGHT,
            "timestamp": TIMESTAMP,
            "tx_hash": "",
            "session_id": item["id"],
            "status": "inactive",
        }
    )
    db["node_statistics"].find_one_and_update(
        {
            "address": item["node_addr"],
            "timestamp": TIMESTAMP.date(),
        },
        {
            "$inc": {
                "session_end_count": 1
            },
        },
    )

# -------------------------------------------------------------------------------------------------------------------- #

db["subscriptions"].update_many(
    {},
    {
        "$rename": {
            "owner": "acc_addr",
            "node": "node_addr",
            "plan": "plan_id",
            "expiry": "inactive_at",
        },
    },
)
db["subscriptions"].update_many(
    {},
    {
        "$unset": {
            "free": 1,
        }
    },
)
db["subscriptions"].update_many(
    {"status": "STATUS_ACTIVE"}, {"$set": {"status": "active"}}
)
db["subscriptions"].update_many(
    {"status": "STATUS_INACTIVE_PENDING"}, {"$set": {"status": "inactive_pending"}}
)
db["subscriptions"].update_many(
    {"status": "STATUS_INACTIVE"}, {"$set": {"status": "inactive"}}
)

cursor = db["subscriptions"].find({"plan_id": 0})
for item in cursor:
    db["subscriptions"].find_one_and_update(
        {
            "id": item["id"],
        },
        {
            "$set": {
                "inactive_at": ZERO_TIMESTAMP
            },
        },
    )

cursor = db["subscriptions"].find({"plan_id": {"$ne": 0}})
for item in cursor:
    plan = db["plans"].find_one(
        {
            "id": item["plan_id"],
        },
    )
    price = [v for v in plan["prices"] if v["denom"] == item["denom"]][0]
    db["subscriptions"].find_one_and_update(
        {
            "id": item["id"],
        },
        {
            "$set": {
                "price": price,
            },
        },
    )

db["subscriptions"].update_many(
    {},
    {
        "$unset": {
            "denom": 1,
        },
    },
)
cursor1 = db["subscriptions"].find({"plan_id": 0, "status": {"$ne": "inactive"}})
for item in cursor1:
    cursor2 = db["sessions"].find(
        {
            "subscription_id": item["id"],
        },
    )
    total = 0
    for session in cursor2:
        if "payment" in session and "amount" in session["payment"]:
            total += int(session["payment"]["amount"])
        if "staking_reward" in session and "amount" in session["staking_reward"]:
            total += int(session["staking_reward"]["amount"])
    refund = {
        "denom": item["deposit"]["denom"],
        "amount": int(item["deposit"]["amount"]) - total,
    }
    db["subscriptions"].find_one_and_update(
        {
            "id": item["id"],
        },
        {
            "$set": {
                "end_height": HEIGHT,
                "end_timestamp": TIMESTAMP,
                "refund": refund,
                "status": "inactive",
                "status_at": HEIGHT,
                "status_timestamp": TIMESTAMP,
            },
        },
    )
    db["events"].insert_one(
        {
            "type": "Subscription.UpdateStatus",
            "height": HEIGHT,
            "timestamp": TIMESTAMP,
            "tx_hash": "",
            "session_id": item["id"],
            "status": "inactive",
        },
    )
    db["node_statistics"].find_one_and_update(
        {
            "address": item["node_addr"],
            "timestamp": TIMESTAMP.date(),
        },
        {
            "$inc": {
                "subscription_end_count": 1
            },
        },
    )

cursor1 = db["subscriptions"].find({"plan_id": {"$ne": 0}})
for item in cursor1:
    plan = db["plans"].find_one(
        {
            "id": item["plan_id"],
        },
    )
    cursor2 = db["subscription_allocations"].find(
        {
            "id": item["id"],
        },
    )
    total = 0
    for alloc in cursor2:
        total += int(alloc["allocated"])
    diff = (plan["gigabytes"] * 1e9) - total

    alloc = db["subscription_allocations"].find_one(
        {
            "id": item["id"],
            "address": item["acc_addr"],
        },
    )
    db["subscription_allocations"].find_one_and_update(
        {
            "id": item["id"],
            "address": item["acc_addr"]
        },
        {
            "$set": {
                "allocated": str(alloc + diff),
            },
        },
    )
    db["events"].insert_one(
        {
            "type": "SubscriptionAllocation.UpdateDetails",
            "height": HEIGHT,
            "timestamp": TIMESTAMP,
            "tx_hash": "",
            "subscription_id": item["id"],
            "acc_addr": item["acc_addr"],
            "granted_bytes": str(alloc + diff),
            "utilised_bytes": alloc["consumed"],
        },
    )

# -------------------------------------------------------------------------------------------------------------------- #

db["subscription_allocations"].update_many(
    {},
    {
        "$rename": {
            "address": "acc_addr",
            "allocated": "granted_bytes",
            "consumed": "utilised_bytes",
        },
    },
)

# -------------------------------------------------------------------------------------------------------------------- #

db["node_statistics"].update_many(
    {},
    {
        "$rename": {
            "address": "addr",
        },
    },
)
