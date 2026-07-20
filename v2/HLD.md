# High-Level Design (HLD)

## System Overview

The Job Scheduling System is designed to handle the submission, scheduling, and execution of one-time and periodic jobs across multiple worker nodes. The system ensures reliable job execution with fault tolerance, minimal duplication, and efficient resource utilization.

---

## Component Details

### 1. Job Submission Service

The Job Submission Service is the entry point for clients to interact with the system. It provides a RESTful API interface for users or services to submit, update, or cancel jobs.

**Responsibilities:**
- Accept job submission requests from clients
- Validate job details and metadata
- Persist job information to the Job Store
- Return a unique Job ID to the client

**API Input Parameters:**
- Job name
- Frequency (One-time, Daily)
- Execution time
- Job payload (task details)

**Processing:**
- Saves job metadata (execution_time, frequency, status = pending) in the Job Store
- Returns a unique Job ID to the client for tracking

---

### 2. Job Store

The Job Store is responsible for persisting job information and maintaining the current state of all jobs and workers in the system. It consists of the following database tables:

#### Job Table
Stores the metadata of each job, including:
- Job ID
- User ID
- Frequency
- Payload
- Execution time
- Retry count
- Status (pending, running, completed, failed)

#### Job Execution Table
Tracks execution attempts for each job. Jobs can be executed multiple times in case of failures.

Stores information including:
- Execution ID
- Start time
- End time
- Worker ID
- Status
- Error message

Each failed job that is retried will have a new entry logged here.

#### Job Schedules Table
Stores scheduling details for each job, including the next_run_time.

**For one-time jobs:**
- next_run_time equals the job's execution time
- last_run_time remains null

**For recurring jobs:**
- next_run_time is updated after each execution to reflect the next scheduled run

#### Worker Table
Stores information about each worker node, including:
- IP address
- Status
- Last heartbeat timestamp
- Capacity
- Current load

---

### 3. Scheduling Service

The Scheduling Service is responsible for selecting jobs for execution based on their next_run_time in the Job Schedules Table.

**Processing Flow:**
1. Periodically queries the Job Schedules Table for jobs scheduled to run at the current minute
2. Retrieves due jobs using a query like: `SELECT * FROM JobSchedulesTable WHERE next_run_time = 1726110000;`
3. Pushes retrieved jobs to the Distributed Job Queue for worker nodes to execute
4. Updates the status in the Job Table to SCHEDULED

---

### 4. Distributed Job Queue

The Distributed Job Queue (e.g., Kafka, RabbitMQ) acts as a buffer between the Scheduling Service and the Execution Service.

**Responsibilities:**
- Hold jobs in a queue structure
- Ensure efficient distribution of jobs to available worker nodes
- Decouple the scheduling and execution layers for better scalability
- Allow the Execution Service to pull jobs and assign them to worker nodes

---

### 5. Execution Service

The Execution Service is responsible for running jobs on worker nodes and updating results in the Job Store. It consists of a Coordinator and a pool of Worker Nodes.

#### Coordinator

The coordinator (or orchestrator) node is responsible for:

**Job Assignment:**
- Distributes jobs from the queue to available worker nodes

**Worker Management:**
- Tracks the status, health, capacity, and workload of active workers

**Failure Handling:**
- Detects when a worker node fails
- Reassigns failed jobs to other healthy nodes

**Load Balancing:**
- Ensures workload is evenly distributed across worker nodes based on available resources and capacity

#### Worker Nodes

Worker nodes are responsible for executing jobs and updating the Job Store with results.

**Execution Flow:**
1. When assigned a job, creates a new entry in the Job Execution Table with status set to running
2. Begins job execution
3. After execution completes, updates the job's final status (completed or failed) along with any output in both the Jobs and Job Execution Table

**Failure Handling:**
- If a worker fails during execution, the coordinator re-queues the job in the Distributed Job Queue
- Another healthy worker can then pick up and complete the job

---

## Data Flow

1. Client submits job via Job Submission Service
2. Job metadata is stored in Job Store
3. Scheduling Service queries Job Schedules Table for due jobs
4. Due jobs are pushed to Distributed Job Queue
5. Coordinator assigns jobs from the queue to available Worker Nodes
6. Worker Nodes execute jobs and update results in Job Store
7. Coordinator monitors worker health and handles failures by reassigning jobs
