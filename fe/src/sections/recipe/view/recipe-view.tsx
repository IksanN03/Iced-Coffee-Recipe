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
import { Inventory } from 'src/sections/inventory/view';

interface Measurement {
  amount: number | null;
  unit: string;
}

interface Recipe {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  sku: string;
  number_of_cups: number | null;
  ingredients: Record<string, Measurement>;
  cogs: number;
}

interface FormData {
  number_of_cups: number | null;
  ingredients: Record<string, Measurement>;
}

const initialFormData: FormData = {
  number_of_cups: null,
  ingredients: {
    'Aren Sugar': { amount: null, unit: 'g' },
    'Milk': { amount: null, unit: 'ml' },
    'Ice Cube': { amount: null, unit: 'g' },
    'Plastic Cup': { amount: null, unit: 'pcs' },
    'Coffee Bean': { amount: null, unit: 'g' },
    'Mineral Water': { amount: null, unit: 'ml' }
  }
};



export function RecipeView() {
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  const [total, setTotal] = useState(0);
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(false);
  const [openModal, setOpenModal] = useState(false);
  const [selectedRecipe, setSelectedRecipe] = useState<Recipe | null>(null);
  const [formData, setFormData] = useState<FormData>(initialFormData);
  const [formErrors, setFormErrors] = useState({
    number_of_cups: '',
    ingredients: {} as Record<string, string>
  });
  const getAvailableUnits = (currentUnit: string) => {
    const unitGroups = {
      weight: ['g', 'kg'],
      volume: ['ml', 'liter'],
      piece: ['pcs']
    };
  
    if (unitGroups.weight.includes(currentUnit.toLowerCase())) {
      return unitGroups.weight;
    }
    if (unitGroups.volume.includes(currentUnit.toLowerCase())) {
      return unitGroups.volume;
    }
    return unitGroups.piece;
  };
  const [toast, setToast] = useState<{
    open: boolean;
    type: 'success' | 'error' | 'warning';
    message: string;
  }>({
    open: false,
    type: 'success',
    message: ''
  });

  const handleCloseToast = () => {
    setToast({ ...toast, open: false });
  };

  const handleCloseModal = () => {
    setFormData({
      number_of_cups: 1,
      ingredients: inventoryItems.reduce((acc: Record<string, Measurement>, item: Inventory) => ({
        ...acc,
        [item.item_name]: { amount: 0, unit: item.uom }
      }), {})
    });
    setSelectedRecipe(null);
    setOpenModal(false);
  };

  const fetchRecipes = useCallback(async () => {
    try {
      setLoading(true);
      const response = await fetch(
        `${import.meta.env.VITE_BASE_URL_BACKEND}/recipe?page=${page + 1}&limit=${limit}&search=${search}`,
        {
          headers: {
            Authorization: `Bearer ${getToken()}`
          }
        }
      );
      const data = await response.json();
      
      if (response.ok) {
        setRecipes(data.data.recipes);
        setTotal(data.data.total_items);
      } else {
        setToast({
          open: true,
          message: data.message.danger || data.message.warning || 'Failed to fetch recipes',
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

  const [inventoryItems, setInventoryItems] = useState<Inventory[]>([]);

const fetchInventory = useCallback(async () => {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_BASE_URL_BACKEND}/inventory`,
      {
        headers: {
          Authorization: `Bearer ${getToken()}`
        }
      }
    );
    const data = await response.json();
    
    if (response.ok) {
      setInventoryItems(data.data.inventory);
      // Create initial ingredients from inventory
      const defaultIngredients = data.data.inventory.reduce((acc: Record<string, Measurement>, item: Inventory) => ({
        ...acc,
        [item.item_name]: { amount: null, unit: item.uom }
      }), {});
      setFormData(prev => ({ ...prev, ingredients: defaultIngredients }));
    
    }
  } catch (error) {
    setToast({
      open: true,
      message: 'Failed to fetch inventory items',
      type: 'error'
    });
  }
}, []);


useEffect(() => {
  fetchInventory();
}, [fetchInventory]);


  useEffect(() => {
    fetchRecipes();
  }, [fetchRecipes]);

  
  // Add validation function
  const validateForm = () => {
    let isValid = true;
    const errors = {
      number_of_cups: '',
      ingredients: {} as Record<string, string>
    };
  
    if (!formData.number_of_cups) {
      errors.number_of_cups = 'Number of cups is required';
      isValid = false;
    }
  

    setFormErrors(errors);
    return isValid;
  };

  const handleSubmit = async () => {
    if (!validateForm()) return;
    try {
      setLoading(true);
      const url = selectedRecipe 
        ? `${import.meta.env.VITE_BASE_URL_BACKEND}/recipe/${selectedRecipe.ID}`
        : `${import.meta.env.VITE_BASE_URL_BACKEND}/recipe`;
      
      const method = selectedRecipe ? 'PUT' : 'POST';

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
        fetchRecipes();
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
        message: 'Failed to save recipe',
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
        onClose={handleCloseToast}
        anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
      >
        <Alert severity={toast.type} onClose={handleCloseToast}>
          {toast.message}
        </Alert>
      </Snackbar>

      <Box display="flex" alignItems="center" mb={5}>
        <Typography variant="h4" flexGrow={1}>
          Recipe Management
        </Typography>
        <TextField
          placeholder="Search recipes..."
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
            setSelectedRecipe(null);
            setOpenModal(true);
          }}
        >
          New Recipe
        </Button>
      </Box>

      <Card>
        <Scrollbar>
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>#</TableCell>
                  <TableCell>SKU</TableCell>
                  <TableCell>Number of Cups</TableCell>
                  <TableCell>Ingredients</TableCell>
                  <TableCell>COGS</TableCell>
                  <TableCell>Created At</TableCell>
                  <TableCell>Updated At</TableCell>
                  <TableCell align="right">Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {recipes.map((recipe, index) => (
                  <TableRow key={recipe.ID}>
                    <TableCell  sx={{ verticalAlign: 'top' }}>{index+1}</TableCell>
                    <TableCell  sx={{ verticalAlign: 'top' }}>{recipe.sku}</TableCell>
                    <TableCell  sx={{ verticalAlign: 'top' }}>{recipe.number_of_cups}</TableCell>
                    <TableCell  sx={{ verticalAlign: 'top' }}>
                    {Object.entries(recipe.ingredients).map(([name, measurement]) => (
                      <Box key={name}>
                        {name}: {measurement.amount} {measurement.unit}
                      </Box>
                    ))}
                  </TableCell>
                    <TableCell  sx={{ verticalAlign: 'top' }}>{recipe.cogs}</TableCell>
                    <TableCell  sx={{ verticalAlign: 'top' }}>{new Date(recipe.CreatedAt).toLocaleString()}</TableCell>
                    <TableCell  sx={{ verticalAlign: 'top' }}>{new Date(recipe.UpdatedAt).toLocaleString()}</TableCell>
                    <TableCell align="right"  sx={{ verticalAlign: 'top' }}>
                      <IconButton
                        onClick={() => {
                          setSelectedRecipe(recipe);
                          setFormData({
                            number_of_cups: recipe.number_of_cups,
                            ingredients: recipe.ingredients
                          });
                          setOpenModal(true);
                        }}
                      >
                        <Iconify icon="eva:edit-fill" />
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
        onClose={handleCloseModal}
      >
        <Box sx={{
          position: 'absolute',
          top: '50%',
          left: '50%',
          transform: 'translate(-50%, -50%)',
          width: 600,
          bgcolor: 'background.paper',
          boxShadow: 24,
          p: 4,
          borderRadius: 1,
          maxHeight: '90vh',
          overflow: 'auto'
        }}>
          <Typography variant="h6" mb={3}>
            {selectedRecipe ? 'Edit Recipe' : 'Add New Recipe'}
          </Typography>
          
          <FormControl fullWidth sx={{ mb: 3 }}>
            <TextField
              label="Number of Cups"
              type="number"
              value={formData.number_of_cups}
              onChange={(e) => setFormData({ 
                ...formData, 
                number_of_cups:parseInt(e.target.value, 10)

              })}
              error={!!formErrors.number_of_cups}
              helperText={formErrors.number_of_cups}
              required
            />
          </FormControl>

          {Object.entries(formData.ingredients).map(([name, measurement]) => (
            <Box key={name} sx={{ mb: 2 }}>
              <Typography variant="subtitle2" mb={1}>
                {name}
              </Typography>
              <Box display="flex" gap={2}>
                <TextField
                  label="Amount"
                  type="number"
                  value={measurement.amount}
                  onChange={(e) => {
                    const newIngredients = { ...formData.ingredients };
                    newIngredients[name] = {
                      ...measurement,
                      amount: parseFloat(e.target.value)
                    };
                    setFormData({ ...formData, ingredients: newIngredients });
                  }}
                  fullWidth
                />
                <FormControl fullWidth>
                  <InputLabel>Unit</InputLabel>
                  <Select
                    value={measurement.unit.toLowerCase()}
                    onChange={(e) => {
                      const newIngredients = { ...formData.ingredients };
                      newIngredients[name] = {
                        ...measurement,
                        unit: e.target.value.toUpperCase()
                      };
                      setFormData({ ...formData, ingredients: newIngredients });
                    }}
                    required
                  >
                     {getAvailableUnits(measurement.unit).map((unit) => (
                      <MenuItem key={unit} value={unit}>
                        {unit}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Box>
            </Box>
          ))}

          <LoadingButton
            fullWidth
            variant="contained"
            onClick={handleSubmit}
            loading={loading}
          >
            {selectedRecipe ? 'Update Recipe' : 'Create Recipe'}
          </LoadingButton>
        </Box>
      </Modal>
    </DashboardContent>
  );
}
