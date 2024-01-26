import math
import sys
from datetime import datetime

from dateutil.parser import parse
from pymongo import MongoClient
from pymongo import UpdateOne, InsertOne

ZERO_TIMESTAMP = parse('0001-01-01T00:00:00.0Z')
HEIGHT = 12_310_005
TIMESTAMP = parse("2023-08-18T12:10:36.572027592Z")
TIMESTAMP_DATE = datetime(year=TIMESTAMP.year, month=TIMESTAMP.month, day=TIMESTAMP.day)

db = MongoClient()["sentinelhub-2"]

cursor, height = db["blocks"].find().sort([("height", -1)]).limit(1), 0
for item in cursor:
    height = item["height"]

print("Latest height", height)
if height != HEIGHT - 1:
    print("Exiting...")
    sys.exit(1)

collections = [
    "deposits",
    "events",
    "nodes",
    "node_statistics"
    "plans",
    "providers",
    "sessions",
    "subscriptions",
    "subscription_quotas"
]
for cname in collections:
    db[cname].drop_indexes()

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
db["events"].update_many({"type": "Plan.AddNode"}, {"$set": {"type": "Plan.LinkNode"}})
db["events"].update_many({"type": "Plan.RemoveNode"}, {"$set": {"type": "Plan.UnlinkNode"}})
db["events"].update_many({"type": "SubscriptionQuota.UpdateDetails"},
                         {"$set": {"type": "SubscriptionAllocation.UpdateDetails"}})
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
            "bandwidth": "internet_speed",
            "handshake": "handshake_dns",
            "price": "gigabyte_prices",
            "reach_status": "health",
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
db["nodes"].update_many({"internet_speed": None}, {"$set": {"internet_speed": {}}})
db["nodes"].update_many({"handshake_dns": None}, {"$set": {"handshake_dns": {}}})
db["nodes"].update_many({"location": None}, {"$set": {"location": {}}})
db["nodes"].update_many({"qos": None}, {"$set": {"qos": {}}})
db["nodes"].update_many({"health": None}, {"$set": {"health": {}}})

nodes_ops = []
cursor = db["nodes"].find()
for item in cursor:
    nodes_ops.append(
        UpdateOne(
            {
                "_id": item["_id"],
            },
            {
                "$set": {
                    "hourly_prices": item["gigabyte_prices"],
                },
            },
        )
    )
db["nodes"].bulk_write(nodes_ops, ordered=False)

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

plans_ops = []
cursor = db["plans"].find()
for item in cursor:
    plans_ops.append(
        UpdateOne(
            {
                "_id": item["_id"],
            },
            {
                "$set": {
                    "gigabytes": int(math.ceil(int(item["gigabytes"]) / 1e9)),
                },
            },
        )
    )
db["plans"].bulk_write(plans_ops, ordered=False)

# -------------------------------------------------------------------------------------------------------------------- #

db["providers"].update_many(
    {},
    {
        "$rename": {
            "address": "addr",
        },
    },
)

providers_ops, events_ops = [], []
cursor = db["providers"].find()
for item in cursor:
    providers_ops.append(
        UpdateOne(
            {
                "_id": item["_id"],
            },
            {
                "$set": {
                    "status": "active",
                    "status_height": HEIGHT,
                    "status_timestamp": TIMESTAMP,
                    "status_tx_hash": "",
                },
            },
        )
    )
    events_ops.append(
        InsertOne(
            {
                "type": "Provider.UpdateDetails",
                "height": HEIGHT,
                "timestamp": TIMESTAMP,
                "tx_hash": "",
                "prov_addr": item["addr"],
                "name": item["name"],
                "identity": item["identity"],
                "website": item["website"],
                "description": item["description"],
                "status": "active",
            }
        )
    )
db["providers"].bulk_write(providers_ops, ordered=False)
db["events"].bulk_write(events_ops, ordered=False)

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

sessions_ops, events_ops, node_statistics_ops = [], [], []
cursor = db["sessions"].find({"status": {"$ne": "inactive"}}).sort([("id", 1)])
for item in cursor:
    print(item["id"])
    sessions_ops.append(
        UpdateOne(
            {
                "_id": item["_id"],
            },
            {
                "$set": {
                    "end_height": HEIGHT,
                    "end_timestamp": TIMESTAMP,
                    "status": "inactive",
                    "status_height": HEIGHT,
                    "status_timestamp": TIMESTAMP,
                },
            },
        )
    )
    events_ops.append(
        InsertOne(
            {
                "type": "Session.UpdateStatus",
                "height": HEIGHT,
                "timestamp": TIMESTAMP,
                "tx_hash": "",
                "session_id": item["id"],
                "status": "inactive",
            },
        )
    )
    node_statistics_ops.append(
        UpdateOne(
            {
                "address": item["node_addr"],
                "timestamp": TIMESTAMP_DATE,
            },
            {
                "$inc": {
                    "session_end_count": 1
                },
            },
            upsert=True,
        )
    )
db["sessions"].bulk_write(sessions_ops, ordered=False)
db["events"].bulk_write(events_ops, ordered=False)
db["node_statistics"].bulk_write(node_statistics_ops, ordered=False)

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

subscriptions_ops = []
cursor = db["subscriptions"].find({"plan_id": {"$ne": 0}}).sort([("id", 1)])
for item in cursor:
    print(item["id"])
    plan = db["plans"].find_one(
        {
            "id": item["plan_id"],
        },
    )
    price = [v for v in plan["prices"] if v["denom"] == item["denom"]][0]
    subscriptions_ops.append(
        UpdateOne(
            {
                "_id": item["_id"],
            },
            {
                "$set": {
                    "price": price,
                },
            },
        )
    )
db["subscriptions"].bulk_write(subscriptions_ops, ordered=False)

db["subscriptions"].update_many(
    {},
    {
        "$unset": {
            "denom": 1,
        },
    },
)

db.sessions.create_index([("subscription_id", 1)])
db.node_statistics.create_index([("address", 1), ("timestamp", 1)])

subscriptions_ops, events_ops, node_statistics_ops = [], [], []
cursor1 = db["subscriptions"].find({"plan_id": 0, "status": {"$ne": "inactive"}}).sort([("id", 1)])
for item in cursor1:
    print(item["id"])
    cursor2 = db["sessions"].find(
        {
            "subscription_id": item["id"],
        },
    )
    total = 0
    for session in cursor2:
        if "payment" in session and \
                session["payment"] is not None and "amount" in session["payment"]:
            total += int(session["payment"]["amount"])
        if "staking_reward" in session and \
                session["staking_reward"] is not None and "amount" in session["staking_reward"]:
            total += int(session["staking_reward"]["amount"])
    refund = {
        "denom": item["deposit"]["denom"],
        "amount": str(int(item["deposit"]["amount"]) - total),
    }
    subscriptions_ops.append(
        UpdateOne(
            {
                "_id": item["_id"],
            },
            {
                "$set": {
                    "end_height": HEIGHT,
                    "end_timestamp": TIMESTAMP,
                    "refund": refund,
                    "status": "inactive",
                    "status_height": HEIGHT,
                    "status_timestamp": TIMESTAMP,
                },
            },
        )
    )
    events_ops.append(
        InsertOne(
            {
                "type": "Subscription.UpdateStatus",
                "height": HEIGHT,
                "timestamp": TIMESTAMP,
                "tx_hash": "",
                "subscription_id": item["id"],
                "status": "inactive",
            },
        )
    )
    node_statistics_ops.append(
        UpdateOne(
            {
                "address": item["node_addr"],
                "timestamp": TIMESTAMP_DATE,
            },
            {
                "$inc": {
                    "subscription_end_count": 1
                },
            },
            upsert=True,
        )
    )
db["subscriptions"].bulk_write(subscriptions_ops, ordered=False)
db["events"].bulk_write(events_ops, ordered=False)
db["node_statistics"].bulk_write(node_statistics_ops, ordered=False)

subscription_allocations_ops, events_ops = [], []
cursor1 = db["subscriptions"].find({"plan_id": {"$ne": 0}}).sort([("id", 1)])
for item in cursor1:
    print(item["id"])
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
    diff = int(plan["gigabytes"] * 1e9) - total
    print(int(plan["gigabytes"] * 1e9), total, diff)
    alloc = db["subscription_allocations"].find_one(
        {
            "id": item["id"],
            "address": item["acc_addr"],
        },
    )
    subscription_allocations_ops.append(
        UpdateOne(
            {
                "id": item["id"],
                "address": item["acc_addr"]
            },
            {
                "$set": {
                    "allocated": str(int(alloc["allocated"]) + diff),
                },
            },
        )
    )
    events_ops.append(
        InsertOne(
            {
                "type": "SubscriptionAllocation.UpdateDetails",
                "height": HEIGHT,
                "timestamp": TIMESTAMP,
                "tx_hash": "",
                "subscription_id": alloc["id"],
                "acc_addr": alloc["address"],
                "granted_bytes": str(int(alloc["allocated"]) + diff),
                "utilised_bytes": alloc["consumed"],
            },
        )
    )
db["subscription_allocations"].bulk_write(subscription_allocations_ops, ordered=False)
db["events"].bulk_write(events_ops, ordered=False)

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

db["node_statistics"].delete_many(
    {
        "address": "",
    },
)
db["node_statistics"].update_many(
    {},
    {
        "$rename": {
            "address": "addr",
        },
    },
)

# -------------------------------------------------------------------------------------------------------------------- #

collections = [
    "deposits",
    "events",
    "nodes",
    "node_statistics"
    "plans",
    "providers",
    "sessions",
    "subscriptions",
    "subscription_quotas"
]
for cname in collections:
    db[cname].drop_indexes()

# -------------------------------------------------------------------------------------------------------------------- #

db["events"].update_many({"type": "Deposit.Add"}, {"$set": {"type": 1}})
db["events"].update_many({"type": "Deposit.Subtract"}, {"$set": {"type": 2}})
db["events"].update_many({"type": "Node.UpdateDetails"}, {"$set": {"type": 3}})
db["events"].update_many({"type": "Node.UpdateStatus"}, {"$set": {"type": 4}})
db["events"].update_many({"type": "Plan.UpdateStatus"}, {"$set": {"type": 5}})
db["events"].update_many({"type": "Plan.LinkNode"}, {"$set": {"type": 6}})
db["events"].update_many({"type": "Plan.UnlinkNode"}, {"$set": {"type": 7}})
db["events"].update_many({"type": "Provider.UpdateDetails"}, {"$set": {"type": 8}})
db["events"].update_many({"type": "Session.UpdateDetails"}, {"$set": {"type": 9}})
db["events"].update_many({"type": "Session.UpdateStatus"}, {"$set": {"type": 10}})
db["events"].update_many({"type": "Subscription.UpdateDetails"}, {"$set": {"type": 11}})
db["events"].update_many({"type": "Subscription.UpdateStatus"}, {"$set": {"type": 12}})
db["events"].update_many({"type": "SubscriptionAllocation.UpdateDetails"}, {"$set": {"type": 13}})
