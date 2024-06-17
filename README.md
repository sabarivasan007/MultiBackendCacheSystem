# Steps for running a spark job
----

### Some Checks in `SparkInit` before deploying the jar

1. Recheck that you are setting `cassandra.privateIp` in `spark.cassandra.connection.host` of `SparkInit` object. 
2. .setMaster() is optional(in case of DSE cluster) in the `SparkInit` object. // Comment it out

**Note:** `SparkInit` is located at `com.precision.etl.utils.spark.`

## Deploying the jar
```bash
sbt clean package
```
* After packaging, jar(without dependencies) will be generated under `target/scala-2.11/` from project's home directory.

```bash
vishnu@dev:~precision_spark_etl$ ls target/scala-2.11/
precision_spark_etl_2.11-0.1.jar  classes  resolution-cache
```

## Copying the jar

* Copy this jar to any of the dse cluster nodes.

## Running the jar
`--master`(Optional) - set the private ip of the master. (Format: dse://private-ip) No port needs to be mentioned.

`--packages` - specify the external dependencies that are used in the project(which are not there in the cluster).

`--class`    - Specify the main class you want to run along with the package.

To confirm that you're giving right path,

* Go into the console by typing `sbt`
* Run `show discoveredMainClasses`.

It will look like something below, copy your class path and specify it in the next step, where you submit the jar to the cluster.

```bash
sbt:precision_spark_etl> show discoveredMainClasses
[info] * com.precision.precision.etl.ingestion.LoadDataIntoSnowflake
[info] * com.precision.precision.etl.ingestion.TableLevelETL
[info] * com.precision.precision.etl.ingestion.experiments.LoadDataSourceIntoSnowflake
[info] * com.precision.precision.etl.staging.CostShares
[info] * com.precision.precision.etl.staging.DimTerritory
[info] * com.precision.precision.etl.staging.Lives
[info] * com.precision.precision.etl.staging.PayerLevelBenType
[info] * com.precision.precision.etl.staging.StagingFormularyDetail
[info] * com.precision.precision.etl.staging.StagingMajorityFormulary
[info] * com.precision.precision.etl.staging.Zip5ToTerr
[info] * com.precision.precision.etl.staging.test
[info] * com.precision.precision.etl.validation.DataValidation
[success] Total time: 7 s, completed 9 Sep, 2019 4:25:05 PM
sbt:precision_spark_etl> 
```

* Finally give the path to the jar in that node you've copied.

For Example, 

```bash
dse spark-submit --packages net.snowflake:spark-snowflake_2.11:2.5.2-spark_2.4,com.typesafe:config:1.3.4,com.amazonaws:aws-java-sdk:1.11.46 --class com.precision.precision.etl.staging.Zip5ToTerr precision_spark_etl_2.11-0.1.jar
```

```bash
dse -u cassandra -p cassandra spark-submit --packages "net.snowflake:spark-snowflake_2.11:2.5.2-spark_2.4,com.typesafe:config:1.3.4,com.amazonaws:aws-java-sdk:1.11.46,com.amazonaws:aws-java-sdk-bom:1.11.682,com.amazonaws:aws-java-sdk-secretsmanager:1.11.682" --executor-memory "7g" --class "com.precision.etl.ingestion.TableLevelETL" /home/ubuntu/urovant_etl_jar/precision_spark_etl_2.11-0.1.jar "IQvia XPT" "TEMP_XPT_PLANTRAK_EXPANDED"

```

* Additionally you can choose submit mode of the job.

`--deploy-mode` 

`cluster`

* Your driver code will run inside any of the nodes inside the cluster.
* No debug info will be printed in the console.
* Choose this mode in case of long running jobs.(For production)

`client`

* Default mode
* Driver code will start from where you are submitting the job.
* Here you could able to see debug informations.
* Choose this when you are testing your code/logic. (Development phase)

---- 

## Loading metadata into casandra 

Download/Export metadata sheet as TSV file and copy into C* table.
```
COPY urovant.data_source_entity(data_vendor,data_name,group_id,entity_name,snowflake_table_name,table_type,ref_table_date_column,has_control_file,filename_affix_type,control_filename_affix_type,filename_affix,control_filename_affix,exclude_filename_affix,table_order,file_type_for_missing_file_validation,file_type,control_file_type,archived_filename_affix,archived_control_filename_affix,schema_type_id,check_for_new_prescriber,"schema",SCHEMA_CDC,schema_type,column_delimiter,row_delimiter,field_enclosed_by,validation_criteria,description,data_arrival_frequency,ingestion_mode,ingestion_type,ingestion_type_auth,ingestion_type_base_uri,s3_bucket,s3_folder,dsefs_path,control_file_dsefs_path,ingestion_type_credentials,table_count,target_date,header,row_duplication_check) FROM 'Urovant_metadata.tsv' WITH HEADER=true AND DELIMITER='@';

```
