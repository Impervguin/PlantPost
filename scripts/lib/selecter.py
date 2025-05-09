from lib.database import PlantDatabase
from typing import List, Dict
import re

class Analyze:
    def __init__(self):
        self.Source = None
        self.PlanningTime = None
        self.ExecutionTime = None
        self.RowsCount = None


class PlantSelecter:
    def __init__(self, dbs: List[PlantDatabase]):
        self.dbs = []
        for db in dbs:
            if not isinstance(db, PlantDatabase):
                raise ValueError('db must be an instance of PlantDatabase')
            self.dbs.append(db)
    
    def AddDatabase(self, db : PlantDatabase):
        if not isinstance(db, PlantDatabase):
            raise ValueError('db must be an instance of PlantDatabase')
        self.dbs.append(db)
    
    def Analyze(self, query) -> Dict[str, Analyze]:
        res = dict()
        if not isinstance(query, str):
            raise ValueError('query must be a string')

        for db in self.dbs:
            analyzeRows = db.AnalyzeQuery(query)
            res[db.ID()] = ParseAnalyze(analyzeRows, db.ID())
        return res

def ParseAnalyze(analyzeRows, source):
    analyze = Analyze()
    analyze.Source = source
    s = " ".join(analyzeRows)
    analyze.PlanningTime = float(re.search(r"[0-9]+\.[0-9]+", re.search(r"Planning Time: [0-9]+\.[0-9]+ ms", s).group(0)).group(0))
    analyze.ExecutionTime = float(re.search(r"[0-9]+\.[0-9]+", re.search(r"Execution Time: [0-9]+\.[0-9]+ ms", s).group(0)).group(0))
    analyze.RowsCount = int(re.search(r"[0-9]+", re.search(r"rows=[0-9]+", re.search(r"\(actual .*rows=[0-9]+ .*\)", s).group(0)).group(0)).group(0))
    return analyze