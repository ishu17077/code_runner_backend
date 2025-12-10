memory_consumer = []
string_size_bytes = 512
num_strings = int((1000000000000000000000 * 1024 * 1024) / string_size_bytes)
memory_consumer = []
for i in range(num_strings):
    memory_consumer.append("A" * string_size_bytes)