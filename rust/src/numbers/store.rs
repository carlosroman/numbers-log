use std::any::Any;

pub trait Store {
    fn save(&mut self, _val: u32) -> Option<bool> {
        None
    }
    fn as_any(&self) -> &dyn Any;
}
