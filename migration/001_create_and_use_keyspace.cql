CREATE KEYSPACE IF NOT EXISTS blog WITH REPLICATION = {'class':'SimpleStrategy', 'replication_factor':1};

USE blog;