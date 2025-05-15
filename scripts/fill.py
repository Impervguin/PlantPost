from lib import database as db
from lib import filler as fl
import os
import argparse
import multiprocessing

class Args:
    batch_size = 1000
    plant_count = 1000
    proc_count = 4
    vervose = False

def ParseArgs():
    parser = argparse.ArgumentParser()
    parser.add_argument('--batch-size', type=int, default=Args.batch_size, dest='batch_size')
    parser.add_argument('--plant-count', type=int, default=Args.plant_count, dest='plant_count', required=True)
    parser.add_argument('--proc-count', type=int, default=Args.proc_count, dest='proc_count')
    parser.add_argument('--verbose', action='store_true', default=Args.vervose, dest='verbose')
    args = parser.parse_args()

    Args.batch_size = args.batch_size
    Args.plant_count = args.plant_count
    Args.proc_count = args.proc_count
    Args.vervose = args.verbose

def poolFunc(fromTo):
    database = db.JsonDatabase(host='localhost', port=5432, database='plantpost', user='impi', password='impi')
    filler = fl.PlantFiller([database])
    fromcnt, tocnt = fromTo
    filler.FillCount(tocnt - fromcnt)
    if Args.vervose:
        print("Filled %d - %d plants" % (fromcnt, tocnt))


if __name__ == '__main__':
    ParseArgs()
    if Args.vervose:
        print("Using batch size: %d" % Args.batch_size)
        print("Using plant count: %d" % Args.plant_count)
        print("Using processes: %d" % Args.proc_count)

    batchCount = Args.plant_count // Args.batch_size
    leftover = Args.plant_count % Args.batch_size
    if Args.vervose:
        print("Filling %d batches of %d plants" % (batchCount, Args.batch_size))
        print("Leftover: %d plants" % leftover)
    packs = [(i*Args.batch_size, (i+1)*Args.batch_size) for i in range(batchCount)]
    if leftover > 0:
        packs.append((batchCount*Args.batch_size, Args.plant_count))
    with multiprocessing.Pool(Args.proc_count) as pool:
        pool.map(poolFunc, packs)
    print("Filling complete")