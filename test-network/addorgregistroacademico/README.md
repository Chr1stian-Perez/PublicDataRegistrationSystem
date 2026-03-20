## Adding Orgregistroacademico to the test network

You can use the `addOrgregistroacademico.sh` script to add another organization to the Fabric test network. The `addOrgregistroacademico.sh` script generates the Orgregistroacademico crypto material, creates an Orgregistroacademico organization definition, and adds Orgregistroacademico to a channel on the test network.

You first need to run `./network.sh up createChannel` in the `test-network` directory before you can run the `addOrgregistroacademico.sh` script.

```
./network.sh up createChannel
cd addOrgregistroacademico
./addOrgregistroacademico.sh up
```

If you used `network.sh` to create a channel other than the default `mychannel`, you need pass that name to the `addorgregistroacademico.sh` script.
```
./network.sh up createChannel -c channel1
cd addOrgregistroacademico
./addOrgregistroacademico.sh up -c channel1
```

You can also re-run the `addOrgregistroacademico.sh` script to add Orgregistroacademico to additional channels.
```
cd ..
./network.sh createChannel -c channel2
cd addOrgregistroacademico
./addOrgregistroacademico.sh up -c channel2
```

For more information, use `./addOrgregistroacademico.sh -h` to see the `addOrgregistroacademico.sh` help text.
