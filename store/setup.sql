-- #FUNCTIONS
-- #1
-- Function to notify the order id 
-- called when an update happens on orders table
CREATE OR REPLACE FUNCTION notify_changes_with_order_id()
RETURNS TRIGGER AS $order_update_event$
BEGIN
    PERFORM pg_notify('updates_to_orders',NEW.id::text);
    RETURN NEW;
END;
$order_update_event$ LANGUAGE plpgsql;

-- #2
-- Function to notify the swap id 
-- called when an update happens on atoomic_swaps table
CREATE OR REPLACE FUNCTION notify_changes_with_swap_id()
RETURNS TRIGGER AS $swap_update_event$
BEGIN
    PERFORM pg_notify('updates_to_atomic_swaps',NEW.id::text);
    RETURN NEW;
END;
$swap_update_event$ LANGUAGE plpgsql;

-- #TRIGGERS
-- #1
-- trigger for insert events on orders
CREATE OR REPLACE TRIGGER notify_orders_insert_trigger
AFTER INSERT ON orders
FOR EACH ROW
EXECUTE FUNCTION notify_changes_with_order_id();

-- #2
-- triger for update events on orders
-- parameter for scope : 1. status
CREATE OR REPLACE TRIGGER notify_orders_update_trigger
AFTER UPDATE ON orders
FOR EACH ROW
WHEN (OLD.status IS DISTINCT FROM NEW.status)
EXECUTE FUNCTION notify_changes_with_order_id();

-- #3
-- triger for update events on atomic_swaps
-- parameter for scope : 1. status
--                       2. current_confirmations
--                       3. filled_amount
CREATE OR REPLACE TRIGGER notify_atomic_swap_trigger
AFTER UPDATE ON atomic_swaps
FOR EACH ROW
WHEN (OLD.status IS DISTINCT FROM NEW.status OR OLD.current_confirmations IS DISTINCT FROM NEW.current_confirmations OR OLD.filled_amount IS DISTINCT FROM NEW.filled_amount)
EXECUTE FUNCTION notify_changes_with_swap_id();