# list
grpcurl  --plaintext 127.0.0.1:50031 list

# describe
grpcurl  --plaintext 127.0.0.1:50031 describe

# add book
grpcurl --plaintext -emit-defaults -d '{"title": "title3","author": "author3"}'\
 127.0.0.1:50031 Bookstore/AddBook

# delete book
grpcurl --plaintext -emit-defaults -d '{"id": 1}'\
 127.0.0.1:50031 Bookstore/DeleteBook

# update book
grpcurl --plaintext -emit-defaults -d '{"id":2, "title": "title2-2","author": "author2-2"}'\
 127.0.0.1:50031 Bookstore/UpdateBook

# list book
grpcurl --plaintext -emit-defaults -d '{}'\
  127.0.0.1:50031 Bookstore/ListBook

# get book
grpcurl --plaintext -emit-defaults -d '{"id":1}'\
  127.0.0.1:50031 Bookstore/GetBook