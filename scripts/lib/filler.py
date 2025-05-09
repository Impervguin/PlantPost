from lib.database import PlantDatabase
from lib.plant import Plant


class PlantFiller:
    def __init__(self, dbs: [PlantDatabase]):
        self.dbs = []
        for db in dbs:
            if not isinstance(db, PlantDatabase):
                raise ValueError('db must be an instance of Database')
            self.dbs.append(db)
    
    def AddDatabase(self, db : PlantDatabase):
        if not isinstance(db, PlantDatabase):
            raise ValueError('db must be an instance of Database')
        self.dbs.append(db)
    
    def FillCount(self, cnt):
        if cnt == 1:
            self.FillOne()
        plants = [Plant.Fake() for _ in range(cnt)]
        for db in self.dbs:
            db.InsertMany(plants)
    
    def FillOne(self):
        plant = Plant.Fake()
        for db in self.dbs:
            db.Insert(plant)


