# Decisions.MD

### Key Technical Decisions and Trade-offs
* **Strict Relational Schemas with Constraints:** Scheduling is naturally relational. An appointment depends on a user existing, and multiple appointments cannot overlap. In SQL databases like **PostgreSQL**, this is naturally handled by the database engine with foreign keys and referential integrity. In contrast, in a NoSQL database, this must be enforced in the application logic because of its schema flexibility.
* **Concurrency:** Part of the requirements is that there can be no overlapping appointments; this is handled in PostgreSQL using **exclusion constraints**, which makes it easier to enforce. With a NoSQL database, one would have to do manual checks to see if a time slot is taken and then do an insert, increasing the risk of race conditions. Using PostgreSQL's exclusion constraints ensures strong data correctness, prevents race conditions, and simplifies application logic. While this might lead to rigid schemas and lower "write" performance, it is better than the high risk of race conditions and complex application logic in NoSQL databases.
* **Component Decomposition (Frontend):** This allowed each component to have its own responsibilities and specific loading/error states for its part of the UI. It results in slightly more code across different files, but it provides significantly higher maintainability and easier debugging.

---

### Assumptions Made About Unclear Requirements
* **Scheduling Conflicts:** The requirement under "Conflict Handling" mentions preventing "scheduling conflicts" without a full definition, while "Concurrent Access" notes multiple users may book simultaneously. I assumed a **Multi-Tenant model**. Overlap checks are globally scoped, preventing both User A and User B from having an appointment at 10:00 AM, and also preventing User A from being double-booked.
* **Time Zone Neutrality:** The requirements didn't state how time zones should be handled, so I assumed an **absolute timing persistence**. If a user in New York books a meeting for 10:00 AM and the service provider is in Lagos, the system automatically handles the time zone offset.
* **Soft Deletes:** It wasn't clear whether to remove records completely or use a soft delete. I implemented **soft deletion** by setting a `deleted_at` timestamp. This ensures auditability if a user claims an appointment disappeared; a hard delete leaves no trace.

---

### What I Would Do Differently With More Time
1.  **Authentication:** Implement a better system using user sessions or tokens to authenticate users before they can use the application.
2.  **Stakeholder Management:** Design the system so stakeholders can manage time slots and turn off bookings for a specific time slot, a specific day, or an entire day.
3.  **Appointment Statuses:** Implement statuses so appointments aren't automatically booked, allowing stakeholders to reject them with clear reasons.
4.  **Websockets:** Currently, a user only finds out a slot is taken when they click "Book." I would implement websockets so the frontend is notified immediately if someone else books a slot while they are looking at it.
5.  **Observability:** Implement proper observability using tools like **Datadog** rather than relying on terminal/console logs.

---

### Open-Ended Considerations

#### Recurring Appointments
* Users can specify that an appointment should be recurring when booking.
* The backend will calculate future dates and run a query to check for overlaps with existing appointments.
* **Success:** Bulk-insert the records.
* **Failure:** Return the specific conflicting date to the user with an error.
* **Durability:** Limit the number of recurrences to prevent the database from becoming bloated.
* **Updates:** Users can choose to update the single record, all future records, or the entire series (past, present, and future).

#### Real-time Updates
To build this, I would use **goroutines** and Postgres `LISTEN/NOTIFY`:
1.  The frontend calls a function like `client.watchAppointments({ userId })`.
2.  The backend starts a loop waiting for a signal from a Go channel.
3.  Postgres triggers an event whenever there is a change to the `appointments` table.
4.  The Go background listener receives the notification and signals that user's channel.
5.  The loop triggers, fetches the latest data, and pushes it to the React app.