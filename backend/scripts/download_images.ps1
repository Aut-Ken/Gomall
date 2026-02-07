$dir = "E:\Gomall\Photos"
if (!(Test-Path -Path $dir)) {
    New-Item -ItemType Directory -Path $dir
}

$images = @(
    @{Url="https://images.unsplash.com/photo-1523275335684-37898b6baf30"; Name="product_watch.jpg"},
    @{Url="https://images.unsplash.com/photo-1505740420928-5e560c06d30e"; Name="product_headphone.jpg"},
    @{Url="https://images.unsplash.com/photo-1526170375885-4d8ecf77b99f"; Name="product_camera.jpg"},
    @{Url="https://images.unsplash.com/photo-1542291026-7eec264c27ff"; Name="product_shoes.jpg"},
    @{Url="https://images.unsplash.com/photo-1583394838336-acd977736f90"; Name="product_ps5.jpg"},
    @{Url="https://images.unsplash.com/photo-1593642632823-8f7856677730"; Name="product_laptop.jpg"},
    @{Url="https://images.unsplash.com/photo-1544947950-fa07a98d237f"; Name="product_book.jpg"},
    @{Url="https://images.unsplash.com/photo-1512496015851-a90fb38ba796"; Name="product_cosmetics.jpg"},
    @{Url="https://images.unsplash.com/photo-1583573636246-18cb2246697f"; Name="product_iphone.jpg"},
    @{Url="https://images.unsplash.com/photo-1560343090-f0409e92791a"; Name="product_leather_shoes.jpg"}
)

foreach ($img in $images) {
    $filename = Join-Path -Path $dir -ChildPath $img.Name
    $url = $img.Url + "?w=800&q=80"
    Write-Host "Downloading $($img.Name)..."
    try {
        Invoke-WebRequest -Uri $url -OutFile $filename -UseBasicParsing
    } catch {
        Write-Host "Failed to download $($img.Name): $_"
    }
}
Write-Host "Download complete."
