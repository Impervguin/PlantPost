import os
import argparse
import psycopg2
import shutil
import sys
from minio import Minio
from minio.error import S3Error

def parse_args():
    parser = argparse.ArgumentParser(description='Sync files between MinIO and filesystem based on database records')
    parser.add_argument('--db-url', required=True, help='PostgreSQL connection URL (e.g. postgresql://user:password@localhost/dbname)')
    parser.add_argument('--minio-endpoint', required=True, help='MinIO server endpoint (e.g. minio.example.com:9000)')
    parser.add_argument('--minio-login', required=True, help='MinIO login (access key)')
    parser.add_argument('--minio-password', required=True, help='MinIO password (secret key)')
    parser.add_argument('--minio-secure', action='store_true', help='Use HTTPS for MinIO connection')
    parser.add_argument('--buckets', required=True, nargs='+', help='List of buckets to sync')
    parser.add_argument('--fs-root', required=True, help='Root directory in filesystem to sync files to')
    parser.add_argument('--min-free-space', type=int, default=1024*1024, help='Minimum free space in bytes')
    return parser.parse_args()

def ensure_fs_directory(path):
    os.makedirs(path, exist_ok=True)


def check_disk_space(required_size, target_path):
    """Check if there's enough free space on the filesystem"""
    stat = shutil.disk_usage(os.path.dirname(target_path))
    return stat.free >= required_size

def get_file_size(minio_client, bucket, file_url):
    """Get file size from MinIO"""
    try:
        obj = minio_client.stat_object(bucket, file_url)
        return obj.size
    except S3Error as e:
        print(f"Error getting file size for {file_url}: {e}", file=sys.stderr)
        return 0

def sync_file_from_minio(minio_client, bucket, file_url, fs_path, min_free_space):
    try:

        file_size = get_file_size(minio_client, bucket, file_url)
        if file_size == 0:
            return False
        
        if not check_disk_space(file_size + min_free_space, fs_path):
            print(f"Not enough disk space for {file_url}. Required: {file_size + min_free_space} bytes", 
                  file=sys.stderr)
            return False

        minio_client.fget_object(bucket, file_url, fs_path)
        print(f"Downloaded {file_url} from bucket {bucket} to {fs_path}")
        return True
    except S3Error as e:
        if e.code == 'NoSuchKey':
            print(f"File {file_url} not found in bucket {bucket}")
        else:
            print(f"Error downloading {file_url}: {e}")
        return False

def sync_file_to_minio(minio_client, bucket, file_url, fs_path):
    try:
        minio_client.fput_object(bucket, file_url, fs_path)
        print(f"Uploaded {fs_path} to bucket {bucket} as {file_url}")
        return True
    except S3Error as e:
        print(f"Error uploading {file_url}: {e}")
        return False

def check_file_in_minio(minio_client, bucket, file_url):
    try:
        minio_client.stat_object(bucket, file_url)
        return True
    except S3Error as e:
        if e.code == 'NoSuchKey':
            return False
        raise

def get_fs_path(fs_root, bucket, file_url):
    """Generate filesystem path in format: <fs_root>/<bucket_name>/<file_url>"""
    file_url = file_url.lstrip('/')
    return os.path.normpath(os.path.join(fs_root, bucket, file_url))

def sync_files(args):
    # Initialize MinIO client
    minio_client = Minio(
        args.minio_endpoint,
        access_key=args.minio_login,
        secret_key=args.minio_password,
        secure=args.minio_secure
    )

    # Connect to PostgreSQL
    conn = psycopg2.connect(args.db_url)
    cur = conn.cursor()

    try:
        # Get all files from database
        cur.execute('SELECT id, name, url FROM "file"')
        files = cur.fetchall()

        for file_id, file_name, file_url in files:
            # Check in which bucket the file exists
            found_in_bucket = None
            for bucket in args.buckets:
                if check_file_in_minio(minio_client, bucket, file_url):
                    found_in_bucket = bucket
                    break

            if found_in_bucket:
                # File exists in MinIO - download to filesystem
                fs_path = get_fs_path(args.fs_root, found_in_bucket, file_url)
                ensure_fs_directory(os.path.dirname(fs_path))
                
                if not os.path.exists(fs_path):
                    sync_file_from_minio(minio_client, found_in_bucket, file_url, fs_path, args.min_free_space)
                else:
                    print(f"File {file_url} already exists in filesystem at {fs_path}")
            else:
                # File doesn't exist in MinIO - check if it exists in any bucket's directory in filesystem
                for bucket in args.buckets:
                    fs_path = get_fs_path(args.fs_root, bucket, file_url)
                    if os.path.exists(fs_path):
                        if sync_file_to_minio(minio_client, bucket, file_url, fs_path):
                            break

    except Exception as e:
        print(f"Error during synchronization: {e}")
    finally:
        cur.close()
        conn.close()

if __name__ == '__main__':
    args = parse_args()
    sync_files(args)