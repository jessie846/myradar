require 'json'
require 'pp'

filename = 'maps/ZJX.json'
file = File.read(filename)
data = JSON.parse(file)
# pp data['facility']['childFacilities'].map { |x| x['id'] }
data['facility'].delete 'childFacilities'
facility = data['facility']
eramConfiguration = facility['eramConfiguration']
nasId = eramConfiguration['nasId']
geoMaps = eramConfiguration['geoMaps']
emergencyChecklist = eramConfiguration['emergencyChecklist']
positionReliefChecklist = eramConfiguration['positionReliefChecklist']
neighboringStarsConfigurations = eramConfiguration['neighboringStarsConfigurations']
eramConfiguration.delete 'asrSites'
eramConfiguration.delete 'beaconCodeBanks'
eramConfiguration.delete 'geoMaps'
pp geoMaps[0]
