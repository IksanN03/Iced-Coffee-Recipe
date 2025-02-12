import { useState, useCallback, useEffect } from 'react';
import {
  Box, Card, Table, Button, Modal, TextField,
  TableBody, Typography, TableContainer, TablePagination,
  FormControl, InputLabel, Select, MenuItem,
  TableHead, TableRow, TableCell, IconButton,
  InputAdornment, Alert, Snackbar
} from '@mui/material';
import { LoadingButton } from '@mui/lab';
import { getToken } from 'src/utils/auth';
import { DashboardContent } from 'src/layouts/dashboard';
import { Iconify } from 'src/components/iconify';
import { Scrollbar } from 'src/components/scrollbar';
import { fCurrency, fNumber } from 'src/utils/format-number';

export interface Inventory {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  item_name: string;
  quantity: number;
  uom: string;
  price_per_qty: number;
}

interface Toast {
  open: boolean;
  message: string;
  type: 'success' | 'error' | 'warning';
}

export function InventoryView() {
  const [inventories, setInventories] = useState<Inventory[]>([]);
  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  const [total, setTotal] = useState(0);
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(false);
  const [openModal, setOpenModal] = useState(false);
  const [selectedItem, setSelectedItem] = useState<Inventory | null>(null);
  const [toast, setToast] = useState<Toast>({
    open: false,
    message: '',
    type: 'success'
  });
  const [formData, setFormData] = useState({
    item_name: '',
    quantity: 0,
    uom: '',
    price_per_qty: 0
  });

  const [formErrors, setFormErrors] = useState({
    item_name: '',
    quantity: '',
    uom: '',
    price_per_qty: ''
  });

  const validateForm = () => {
    let isValid = true;
    const errors = {
      item_name: '',
      quantity: '',
      uom: '',
      price_per_qty: ''
    };

    if (!formData.item_name.trim()) {
      errors.item_name = 'Item name is required';
      isValid = false;
    }

    if (formData.quantity <= 0) {
      errors.quantity = 'Quantity must be greater than 0';
      isValid = false;
    }

    if (!formData.uom) {
      errors.uom = 'UOM is required';
      isValid = false;
    }

    if (formData.price_per_qty <= 0) {
      errors.price_per_qty = 'Price must be greater than 0';
      isValid = false;
    }

    setFormErrors(errors);
    return isValid;
  };

  const fetchInventories = useCallback(async () => {
    try {
      setLoading(true);
      const response = await fetch(
        `${import.meta.env.VITE_BASE_URL_BACKEND}/inventory?page=${page + 1}&limit=${limit}&search=${search}`,
        {
          headers: {
            Authorization: `Bearer ${getToken()}`
          }
        }
      );
      const data = await response.json();
      
      if (response.ok) {
        setInventories(data.data.inventory);
        setTotal(data.data.total_items);
      } else {
        setToast({
          open: true,
          message: data.message.danger || data.message.warning || 'Failed to fetch inventory',
          type: data.message.danger ? 'error' : 'warning'
        });
      }
    } catch (error) {
      setToast({
        open: true,
        message: 'Network error occurred',
        type: 'error'
      });
    } finally {
      setLoading(false);
    }
  }, [page, limit, search]);

  useEffect(() => {
    fetchInventories();
  }, [fetchInventories]);

  const handleSubmit = async () => {
    if (!validateForm()) return;

    try {
      setLoading(true);
      const url = selectedItem 
        ? `${import.meta.env.VITE_BASE_URL_BACKEND}/inventory/${selectedItem.ID}`
        : `${import.meta.env.VITE_BASE_URL_BACKEND}/inventory`;
      
      const method = selectedItem ? 'PUT' : 'POST';

      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${getToken()}`
        },
        body: JSON.stringify(formData)
      });

      const data = await response.json();

      if (response.ok) {
        setToast({
          open: true,
          message: data.message.success,
          type: 'success'
        });
        setOpenModal(false);
        fetchInventories();
      } else {
        setToast({
          open: true,
          message: data.message.danger || data.message.warning,
          type: data.message.danger ? 'error' : 'warning'
        });
      }
    } catch (error) {
      setToast({
        open: true,
        message: 'Failed to save inventory item',
        type: 'error'
      });
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm('Are you sure you want to delete this item?')) return;

    try {
      setLoading(true);
      const response = await fetch(
        `${import.meta.env.VITE_BASE_URL_BACKEND}/inventory/${id}`, 
        {
          method: 'DELETE',
          headers: {
            Authorization: `Bearer ${getToken()}`
          }
        }
      );

      const data = await response.json();

      if (response.ok) {
        setToast({
          open: true,
          message: data.message.success,
          type: 'success'
        });
        fetchInventories();
      } else {
        setToast({
          open: true,
          message: data.message.danger || data.message.warning,
          type: data.message.danger ? 'error' : 'warning'
        });
      }
    } catch (error) {
      setToast({
        open: true,
        message: 'Failed to delete inventory item',
        type: 'error'
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <DashboardContent>
      <Snackbar
        open={toast.open}
        autoHideDuration={6000}
        onClose={() => setToast({ ...toast, open: false })}
        anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
      >
        <Alert severity={toast.type} onClose={() => setToast({ ...toast, open: false })}>
          {toast.message}
        </Alert>
      </Snackbar>

      <Box display="flex" alignItems="center" mb={5}>
        <Typography variant="h4" flexGrow={1}>
          Inventory Management
        </Typography>
        <TextField
          placeholder="Search inventory..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          sx={{ mr: 2 }}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <Iconify icon="eva:search-fill" />
              </InputAdornment>
            ),
          }}
        />
        <Button
          variant="contained"
          startIcon={<Iconify icon="eva:plus-fill" />}
          onClick={() => {
            setSelectedItem(null);
            setFormData({
              item_name: '',
              quantity: 0,
              uom: '',
              price_per_qty: 0
            });
            setFormErrors({
              item_name: '',
              quantity: '',
              uom: '',
              price_per_qty: ''
            });
            setOpenModal(true);
          }}
        >
          New Item
        </Button>
      </Box>

      <Card>
        <Scrollbar>
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>#</TableCell>
                  <TableCell>Item Name</TableCell>
                  <TableCell>Quantity</TableCell>
                  <TableCell>UOM</TableCell>
                  <TableCell>Price Per Qty</TableCell>
                  <TableCell>Created At</TableCell>
                  <TableCell>Updated At</TableCell>
                  <TableCell align="right">Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {inventories.map((item, index) => (
                  <TableRow key={item.ID}>
                    <TableCell>{index+1}</TableCell>
                    <TableCell>{item.item_name}</TableCell>
                    <TableCell>{fNumber(item.quantity)}</TableCell>
                    <TableCell>{item.uom}</TableCell>
                    <TableCell>{ fCurrency(item.price_per_qty)}</TableCell>
                    <TableCell>{new Date(item.CreatedAt).toLocaleString()}</TableCell>
                    <TableCell>{new Date(item.UpdatedAt).toLocaleString()}</TableCell>
                    <TableCell align="right">
                      <IconButton
                        onClick={() => {
                          setSelectedItem(item);
                          setFormData(item);
                          setOpenModal(true);
                        }}
                      >
                        <Iconify icon="eva:edit-fill" />
                      </IconButton>
                      <IconButton 
                        onClick={() => handleDelete(item.ID)}
                        color="error"
                      >
                        <Iconify icon="eva:trash-2-fill" />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </Scrollbar>

        <TablePagination
          page={page}
          component="div"
          count={total}
          rowsPerPage={limit}
          onPageChange={(e, newPage) => setPage(newPage)}
          onRowsPerPageChange={(e) => {
            setLimit(parseInt(e.target.value, 10));
            setPage(0);
          }}
        />
      </Card>

      <Modal
        open={openModal}
        onClose={() => setOpenModal(false)}
      >
        <Box sx={{
          position: 'absolute',
          top: '50%',
          left: '50%',
          transform: 'translate(-50%, -50%)',
          width: 400,
          bgcolor: 'background.paper',
          boxShadow: 24,
          p: 4,
          borderRadius: 1
        }}>
          <Typography variant="h6" mb={3}>
            {selectedItem ? 'Edit Inventory Item' : 'Add New Inventory Item'}
          </Typography>
          
          <FormControl fullWidth sx={{ mb: 2 }}>
            <TextField
              label="Item Name"
              value={formData.item_name}
              onChange={(e) => setFormData({ ...formData, item_name: e.target.value })}
              error={!!formErrors.item_name}
              helperText={formErrors.item_name}
              required
            />
          </FormControl>

          <FormControl fullWidth sx={{ mb: 2 }}>
          <TextField
            label="Quantity"
            type="text"
            value={formData.quantity ? fNumber(formData.quantity) : ''}
            onKeyPress={(e) => {
              if (!/[0-9]/.test(e.key)) {
                e.preventDefault();
              }
            }}
            onChange={(e) => {
              const rawValue = e.target.value.replace(/\D/g, '');
              setFormData({ 
                ...formData, 
                quantity: rawValue ? Number(rawValue) : 0 
              });
            }}
            error={!!formErrors.quantity}
            helperText={formErrors.quantity}
            required
          />
        </FormControl>

          <FormControl fullWidth sx={{ mb: 2 }}>
            <InputLabel>UOM</InputLabel>
            <Select
             value={formData.uom.toLowerCase()}
             onChange={(e) => setFormData({ ...formData, uom: e.target.value.toLowerCase() })}
             error={!!formErrors.uom}
              required
            >
              <MenuItem value="pcs">pcs</MenuItem>
              <MenuItem value="kg">kg</MenuItem>
              <MenuItem value="liter">liter</MenuItem>
              <MenuItem value="ml">ml</MenuItem>
              <MenuItem value="g">g</MenuItem>
            </Select>
            {formErrors.uom && (
              <Typography color="error" variant="caption">
                {formErrors.uom}
              </Typography>
            )}
          </FormControl>

          <FormControl fullWidth sx={{ mb: 3 }}>
          <TextField
            label="Price Per Quantity"
            type="text"
            value={formData.price_per_qty ? fCurrency(formData.price_per_qty) : ''}
            onKeyPress={(e) => {
              if (!/[0-9]/.test(e.key)) {
                e.preventDefault();
              }
            }}
            onChange={(e) => {
              const rawValue = e.target.value.replace(/\D/g, '');
              setFormData({ 
                ...formData, 
                price_per_qty: rawValue ? Number(rawValue) : 0 
              });
            }}
            error={!!formErrors.price_per_qty}
            helperText={formErrors.price_per_qty}
            required
          />
        </FormControl>

          <LoadingButton
            fullWidth
            variant="contained"
            onClick={handleSubmit}
            loading={loading}
          >
            {selectedItem ? 'Update Item' : 'Create Item'}
          </LoadingButton>
        </Box>
      </Modal>
    </DashboardContent>
  );
}
