package database

/*

db.sync_statuses.createIndex({'app_name': 1}, {unique: true})
db.blocks.createIndex({'height': 1}, {unique: true})
db.txs.createIndex({'height': 1, 'result.code': 1}, {partialFilterExpression: {'result.code': 0}})

db.nodes.createIndex({'address': 1}, {unique: true})
db.deposits.createIndex({'address': 1}, {unique: true})
db.subscriptions.createIndex({'id': 1}, {unique: true})
db.allocations.createIndex({'id': 1, 'address': 1}, {unique: true})
db.sessions.createIndex({'id': 1}, {unique: true})
db.plans.createIndex({'id': 1}, {unique: true})
db.providers.createIndex({'address': 1}, {unique: true})

db.nodes.createIndex({'remote_url': 1, 'status': 1}, {partialFilterExpression: {'remote_url': {$exists: true}}})

*/
