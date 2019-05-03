import os
import subprocess
import codecs
import pyAesCrypt
import base64

# Take input of file name
filename = input("Enter file name : ")
# encryption/decryption buffer size - 64K
bufferSize = 64 * 1024
# file information
statinfo = os.stat(filename)
# n is the number of fragments 
n = 35
# size of each fragment
frag_size = int(statinfo.st_size/n)
# size of last fragment
last_frag_size = statinfo.st_size%frag_size
# command to split the file into fragments
cmd = "split -d -a2 --bytes="+str(frag_size) + " " + filename + " file"
os.system(cmd)
# randomly generated secret key used for encryption and decryption
secret_key = str(os.urandom(32))
decrypted_files = []
for i in range(36):
    if(i<10):
        s = str(i).zfill(2)
    else:
        s = str(i)
    s = "file" + s
    encrypted_file = "encrypted" + s
    decrypted_file = "decrypted" + s
    decrypted_files.append(decrypted_file)
    # encrypt each file into new files with names "encrypted00","encrypted01",etc
    pyAesCrypt.encryptFile(s, encrypted_file, secret_key, bufferSize)
    # decrypt each file into new files with names "decrypted00","decrypted01",etc
    pyAesCrypt.decryptFile(encrypted_file, decrypted_file, secret_key, bufferSize)

# append all decrypted files into one final file (should be same as the original file)
with open('final_file','wb') as outfile:
    for fname in decrypted_files:
        with open(fname,'rb') as infile:
            outfile.write(infile.read())

