import os
import psycopg2
import sys
import pandas as pd
import requests
import csv
import uuid
import time
from database import openDataBaseConnection, closeDataBaseConnection, closeDataBaseConnectionWithoutCommit

# Check env
env = "prod"

DCM_X_API_KEY = '35a8e7ce4385450bb35c50e1e9f1dcfd'

SAAMS_TOKEN = 'eyJraWQiOiJ4R2pRdVJWeVM1MVZKTEJsUXhIY3BheWpBWFwvdFwvOTlGQytQdzlVUlpGVEk9IiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiIwNGExZmMzYy04OTFhLTRkODItYWQyZS1hMzNkNmZiODM1NDMiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmFwLXNvdXRoLTEuYW1hem9uYXdzLmNvbVwvYXAtc291dGgtMV9GSnZYQjRLUGoiLCJwaG9uZV9udW1iZXJfdmVyaWZpZWQiOnRydWUsImNvZ25pdG86dXNlcm5hbWUiOiIxM2UzZTNmMy04YjRkLTQwYjMtOTIzYy1jNGQ0MTI3MzBjNTIiLCJnaXZlbl9uYW1lIjoiU2hyYW1lZXQgTmF5YWsiLCJhdWQiOiIyNGxxYzg3bmpkczVhcWprbmlib2hyNWVybyIsImV2ZW50X2lkIjoiNWJiNTQ3NjUtYTg1NC00MTlkLTk4MTQtNjU2ZWJkNGMyZWI1IiwidG9rZW5fdXNlIjoiaWQiLCJhdXRoX3RpbWUiOjE3MjUwODMzODYsInBob25lX251bWJlciI6Iis5MTcwMzg1NDczMjAiLCJleHAiOjE3MjUxMTY3NTcsImlhdCI6MTcyNTExMzE1NywiZW1haWwiOiJuYXlha3NocmFtZWV0MzNAZ21haWwuY29tIn0.dCECab-L6mgsOUeLm2f6hDU7YLM6SKVIuRJuMP6a9_22Hi8TLmphq9lWOT4u_3jBn_aqGbvlsJ_OhKi5nTe1Lysk_m-20cYb-Ri57droVUDwe5n4pXyNW49X6yqbNqAqxr6TSxhhddiO5_ehUr3WoWCdr87IvxXf0A-2nRXZhW5TIwdtngQcl2EiMLvhqfMD-7TfHuAgb-t9UjsxtLGuiI8BSnefhXB6E724hc8PZQwdi0ykgF0ZX17O-Mn_VKNXIxWx74B-ILMnypqCBZEYvCu9zV6xfSb07aQEI5-H-KJ1puZaUAP0phuEP96w_KKBrY36mjQ6WclZLJGYHmNqmQ'

def getDevicesFromAccessPointID(accessPointId):

    try:
        cursor.execute('''select d."serialNumber" , d."organisationId"
            from devices d 
            inner join access_point_devices apd 
            on apd."serialNumber" = d."serialNumber" 
            inner join access_points ap 
            on ap."accessPointId" = apd ."accessPointId" 
            where ap."accessPointId" = %s;''',(accessPointId,))
        devices = cursor.fetchall()
    
        return devices

    except Exception as e:
            print(e)

def resetConfigDevice(orgId, serialNumber, status):

    try:
        if status == 'reset':

            url ='https://iot.api.spintly.com/deviceConfigurationManagement/internal/v1/devices/'+ serialNumber +'/reset'

        elif status == 'configure':

            url ='https://iot.api.spintly.com/deviceConfigurationManagement/internal/v1/devices/'+ serialNumber +'/configure'

        else:
            print("resetConfigDevice: status is wrong")
            return
        
        headers = {'Accept': '*/*','Content-Type':'application/json','X-API-KEY':DCM_X_API_KEY}

        response = requests.post(url=url,headers=headers)
        print(response)
        print(response.json())

        return response

    except Exception as e:
            print(e)

def checkDeviceStatus(serialNumber ,status):

    try:

        headers = {'Accept': '*/*','Content-Type':'application/json','X-API-KEY':DCM_X_API_KEY}

        url ='https://iot.api.spintly.com/deviceConfigurationManagement/internal/v1/getConfigurationStatus/device/' + serialNumber

        response = requests.get(url=url,headers=headers)

        resMessage = response.json()
        print(resMessage)

        ConfigurationState = resMessage['message']['ConfigurationState']
        
        if ConfigurationState == status:
            return True
        else:         
            return False

    except Exception as e:
            print(e)

def permissionToChange(orgId, accessPointId, usersWithPermissions, action):

    try:

        headers = {'Accept': '*/*','Content-Type':'application/json','Authorization':SAAMS_TOKEN}

        url ='https://saams.api.spintly.com/accessManagementV3/v1/organisations/'+ str(orgId) +'/accessPoint/'+ str(accessPointId) +'/users/permissions'

        if action == 'remove':

            payload = {"permissionsToAdd":[],"permissionsToRemove":usersWithPermissions,"pendingPermissionsToRemove":[]}

        elif action == 'add':

            payload = {"permissionsToAdd":usersWithPermissions,"permissionsToRemove":[],"pendingPermissionsToRemove":[]}

        else:
            print("Wrong Action")

            return
        
        print(url, payload, headers)
        response = requests.patch(url=url,headers=headers, json= payload)
        print(response)
        print(response.json())

        return response
    
    except Exception as e:
        print(e)


try:    
    accessPointId = '9467'
    orgId = '1036'

    print("AccessPointId is ",accessPointId)
    fileOutputName = "output_" + accessPointId + ".txt"

    output_file = open(fileOutputName, 'w', encoding='utf-8')
    original_stdout = sys.stdout
    sys.stdout = output_file

    connection,cursor = openDataBaseConnection(env)

    try:

        # // Select accessPointId and orgId//
       

        headers = {'Accept': '*/*','Content-Type':'application/json','Authorization':SAAMS_TOKEN}

        url = "https://saams.api.spintly.com/accessManagementV3/v1/organisations/"+ orgId +"/accessPoint/"+ accessPointId +"/users/permissions"

        print("URL:", url,"HEADER:", headers)

        response = requests.get(url=url,headers=headers)

        resStatusCode = response.status_code        

        if resStatusCode == 200:
            resMessage = response.json()

            print(resMessage)
            
            usersWithPermissions = resMessage['message']['usersWithPermissions']

            usersWithPermissions_csv_file = 'usersWithPermissions_{}.csv'.format(accessPointId) 

            with open(usersWithPermissions_csv_file, mode='w', newline='') as file:
                writer = csv.writer(file)
                writer.writerow(['usersWithPermissions'])

            for user in usersWithPermissions:
                with open(usersWithPermissions_csv_file, mode='a', newline='') as file:
                    writer = csv.writer(file)
                    writer.writerow([user['id']])

            df = pd.read_csv(usersWithPermissions_csv_file)

            usersWithPermissions = df['usersWithPermissions'].tolist()

            resp = permissionToChange(orgId, accessPointId, usersWithPermissions, 'remove')
            if resp.status_code == 200:
                print("Sucessifully removed perms")
            else:
                print("Could not delete perms. exiting...")
                exit(1)


            usersWithPermissions_csv_file = 'usersWithPermissions_{}.csv'.format(accessPointId) 
            df = pd.read_csv(usersWithPermissions_csv_file)
            usersWithPermissions = df['usersWithPermissions'].tolist()

            devices = getDevicesFromAccessPointID(accessPointId)
            
            deviceUnConfigure = False
            
            deviceConfigure = False 

            if devices != []:
                for device in devices:
                    serialNumber = device[0]
                    orgId = device[1]
                    
                    response = resetConfigDevice(orgId, serialNumber, 'reset')

                    resStatusCode = response.status_code        

                    if resStatusCode == 200:
                        print("Sucessfuly resetted device: ", serialNumber) 
                    else:
                        print("coudl not reset device", serialNumber)
                        exit(2)  

                temp = True
                

                while deviceUnConfigure != True:
                    time.sleep(10)
                    for device in devices:
                        serialNumber = device[0]
                        orgId = device[1]

                        if checkDeviceStatus(serialNumber, 'unconfigured'):
                            if temp == True:
                                deviceUnConfigure = True
                                temp = True
                            else:
                                deviceUnConfigure = False
                                temp = False

                approval = 'n'

                while approval != 'y':
                    approval = input( "Do you want to cont. with configuring anf addding permissions? (yes/no)")

                if deviceUnConfigure:
                    for device in devices:
                        serialNumber = device[0]
                        orgId = device[1]
                        response = resetConfigDevice(orgId, serialNumber, 'configure')

                        if response.status_code == 200:
                            print("Confgured device ", serialNumber)
                        else:
                            print("failed to configure device", serialNumber)     
                            exit(3)


                    print("Sleeping for 1 min")
                    time.sleep(10)

                

                response = permissionToChange(orgId, accessPointId, usersWithPermissions, 'add')
                if response.status_code == 200:
                    print("Successfully added permissions in orgID " + orgId + " on accessPoint " +accessPointId)   
                else:
                    print("failed to add permissions in orgID " + orgId + " on accessPoint " +accessPointId) 
                    exit(4)  
            else:
                print("no devices found for AP")
                exit(6)

        else:
            print("User permission api failed with error : ",response.json())
            exit(5)
                

    except Exception as e:
        print(e)

        sys.stdout = original_stdout

    # Close the file
        output_file.close()
        closeDataBaseConnectionWithoutCommit(connection,cursor)
        

except Exception as e:
        print(e)
        